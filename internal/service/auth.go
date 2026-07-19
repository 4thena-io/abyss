package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"github.com/4thena-io/abyss/internal/auth"
	"github.com/4thena-io/abyss/internal/integration/forge"
	"github.com/4thena-io/abyss/internal/model"
	"golang.org/x/oauth2"
)

type UserRepository interface {
	Save(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetAll(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByForgeID(ctx context.Context, forgeID int64) (*model.User, error)
	GetByTokenLastEight(ctx context.Context, lastEight string) ([]model.User, error)
	Count(ctx context.Context) (int64, error)
}

type AuthService struct {
	userRepository UserRepository
	forge          forge.Forge
	clientID       string
	clientSecret   string
	forgeType      string
	forgeHost      string
	callbackURL    string
	sessionSecret  string
	owner          string
}

func NewAuthService(
	userRepository UserRepository,
	forge forge.Forge,
	clientID, clientSecret, forgeType, forgeHost, callbackURL, sessionSecret, owner string,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		forge:          forge,
		clientID:       clientID,
		clientSecret:   clientSecret,
		forgeType:      forgeType,
		forgeHost:      forgeHost,
		callbackURL:    callbackURL,
		sessionSecret:  sessionSecret,
		owner:          owner,
	}
}

func (s *AuthService) ForgeType() string { return s.forgeType }
func (s *AuthService) ForgeHost() string { return s.forgeHost }

func (s *AuthService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.userRepository.GetAll(ctx)
}

func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.userRepository.GetByUsername(ctx, username)
}

// oauthConfig builds the OAuth2 config for the configured forge, delegating
// the endpoint paths to the forge implementation since each forge follows
// the same OAuth2 spec but exposes it at different paths.
func (s *AuthService) oauthConfig() *oauth2.Config {
	authPath, tokenPath := "/login/oauth/authorize", "/login/oauth/access_token"
	if s.forge != nil {
		authPath, tokenPath = s.forge.OAuthEndpoints()
	}
	return &oauth2.Config{
		ClientID:     s.clientID,
		ClientSecret: s.clientSecret,
		RedirectURL:  s.callbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  s.forgeHost + authPath,
			TokenURL: s.forgeHost + tokenPath,
		},
	}
}

// AuthURL returns the URL to redirect the user to for OAuth2 authorization.
// The state is a random value used to prevent CSRF attacks — you generate it,
// store it temporarily (cookie/session), and verify it matches in the callback.
func (s *AuthService) AuthURL(state string) string {
	return s.oauthConfig().AuthCodeURL(state)
}

// Exchange trades the authorization code (received in the callback) for an
// access token. This is a server-to-server request — the client_secret never
// leaves the server. The code is single-use and short-lived.
func (s *AuthService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.oauthConfig().Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

// GetForgeUser uses the user's OAuth access token to fetch their identity from
// the forge. This is separate from GetAuthenticatedUser which uses the bot token.
func (s *AuthService) GetForgeUser(ctx context.Context, token *oauth2.Token) (*model.ForgeUser, error) {
	if s.forge == nil {
		return nil, fmt.Errorf("forge is unavailable")
	}
	return s.forge.GetUserByToken(ctx, token.AccessToken)
}

// GetOrCreateUser finds an existing user by their forge ID, or creates a new one.
// The very first user to log in becomes admin — no selection needed.
func (s *AuthService) GetOrCreateUser(ctx context.Context, forgeUser *model.ForgeUser) (*model.User, error) {
	if s.forge == nil {
		return nil, fmt.Errorf("forge is unavailable")
	}
	member, err := s.forge.IsMemberOfOwner(ctx, s.owner, forgeUser.Username)
	if err != nil {
		return nil, err
	}
	if !member {
		return nil, ErrUnauthorized
	}

	user, err := s.userRepository.GetByForgeID(ctx, forgeUser.ID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		user.AvatarURL = forgeUser.AvatarURL
		user.Email = forgeUser.Email
		if err := s.userRepository.Update(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// First user ever = admin
	count, err := s.userRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	user = &model.User{
		ForgeID:   forgeUser.ID,
		Username:  forgeUser.Username,
		Email:     forgeUser.Email,
		IsAdmin:   count == 0,
		AvatarURL: forgeUser.AvatarURL,
	}
	if err := s.userRepository.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// GenerateJWT signs a short-lived JWT for browser sessions (cookie).
func (s *AuthService) GenerateJWT(user *model.User) (string, error) {
	return auth.Sign(user.ID, user.Username, user.IsAdmin, user.AvatarURL, s.sessionSecret)
}

// ValidateJWT verifies a JWT and returns its claims.
func (s *AuthService) ValidateJWT(tokenStr string) (*auth.Claims, error) {
	return auth.Verify(tokenStr, s.sessionSecret)
}

// tokenLastEightLen is how many trailing characters of a PAT are kept in the
// clear to make lookup an indexed query instead of a full-table hash compare.
const tokenLastEightLen = 8

// hashToken returns the hex-encoded SHA-256 digest of a PAT. The token itself
// is generated from 32 bytes of crypto/rand (256 bits of entropy), so it
// already carries enough randomness to stand in as its own salt.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// GeneratePersonalToken creates a random PAT, persists only its hash, and
// returns the raw value. Used by the CLI — the raw token is shown once here
// and never stored, so a DB leak alone can't be used to authenticate.
func (s *AuthService) GeneratePersonalToken(ctx context.Context, user *model.User) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(b)

	hash := hashToken(token)
	user.TokenHash = &hash
	user.TokenLastEight = token[len(token)-tokenLastEightLen:]
	if err := s.userRepository.Update(ctx, user); err != nil {
		return "", err
	}
	return token, nil
}

// GetUserByID fetches a user by their primary key.
func (s *AuthService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepository.GetByID(ctx, id)
}

// GetUserByToken looks up a user by their PAT. Used by the auth middleware.
// The last-eight-character index narrows candidates to a handful of rows,
// then a constant-time compare of the full hash confirms the match.
func (s *AuthService) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	if len(token) < tokenLastEightLen {
		return nil, nil
	}

	candidates, err := s.userRepository.GetByTokenLastEight(ctx, token[len(token)-tokenLastEightLen:])
	if err != nil {
		return nil, err
	}

	hash := []byte(hashToken(token))
	for i := range candidates {
		u := candidates[i]
		if u.TokenHash != nil && subtle.ConstantTimeCompare([]byte(*u.TokenHash), hash) == 1 {
			return &u, nil
		}
	}
	return nil, nil
}

// RevokePersonalToken clears the user's PAT.
func (s *AuthService) RevokePersonalToken(ctx context.Context, user *model.User) error {
	user.TokenHash = nil
	user.TokenLastEight = ""
	return s.userRepository.Update(ctx, user)
}
