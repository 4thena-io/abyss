package service

import (
	"context"
	"crypto/rand"
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
	GetByForgeID(ctx context.Context, forgeID int64) (*model.User, error)
	GetByToken(ctx context.Context, token string) (*model.User, error)
	Count(ctx context.Context) (int64, error)
}

type AuthService struct {
	userRepository UserRepository
	forge          forge.Forge
	clientID       string
	clientSecret   string
	forgeHost      string
	callbackURL    string
	sessionSecret  string
	owner          string
}

func NewAuthService(
	userRepository UserRepository,
	forge forge.Forge,
	clientID, clientSecret, forgeHost, callbackURL, sessionSecret, owner string,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		forge:          forge,
		clientID:       clientID,
		clientSecret:   clientSecret,
		forgeHost:      forgeHost,
		callbackURL:    callbackURL,
		sessionSecret:  sessionSecret,
		owner:          owner,
	}
}

// oauthConfig builds the OAuth2 config for the configured forge.
// Each forge follows the same OAuth2 spec but has different endpoint paths.
func (s *AuthService) oauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.clientID,
		ClientSecret: s.clientSecret,
		RedirectURL:  s.callbackURL,
		Endpoint: oauth2.Endpoint{
			// Gitea and Forgejo use these paths. GitHub uses different ones.
			AuthURL:  s.forgeHost + "/login/oauth/authorize",
			TokenURL: s.forgeHost + "/login/oauth/access_token",
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
	return s.forge.GetUserByToken(ctx, token.AccessToken)
}

// GetOrCreateUser finds an existing user by their forge ID, or creates a new one.
// The very first user to log in becomes admin — no selection needed.
func (s *AuthService) GetOrCreateUser(ctx context.Context, forgeUser *model.ForgeUser) (*model.User, error) {
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

// GeneratePersonalToken creates a random PAT, persists it, and returns the raw value.
// Used by the CLI — the token is shown once and stored hashed or plain depending on policy.
func (s *AuthService) GeneratePersonalToken(ctx context.Context, user *model.User) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(b)

	user.Token = &token
	if err := s.userRepository.Update(ctx, user); err != nil {
		return "", err
	}
	return token, nil
}

// GetUserByToken looks up a user by their PAT. Used by the auth middleware.
func (s *AuthService) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	return s.userRepository.GetByToken(ctx, token)
}
