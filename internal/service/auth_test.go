package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/4thena-io/abyss/internal/integration/forge/mocks"
	"github.com/4thena-io/abyss/internal/model"
	svcmocks "github.com/4thena-io/abyss/internal/service/mocks"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

func TestAuthService_AuthURL(t *testing.T) {
	svc := NewAuthService(svcmocks.NewMockUserRepository(t), nil, "client-id", "client-secret", "gitea", "https://forge.example.com", "https://abyss.example.com/callback", "secret", "owner")

	url := svc.AuthURL("random-state")
	if !strings.Contains(url, "client_id=client-id") || !strings.Contains(url, "state=random-state") {
		t.Fatalf("unexpected auth URL: %s", url)
	}
	if !strings.HasPrefix(url, "https://forge.example.com/login/oauth/authorize") {
		t.Fatalf("expected authorize endpoint on the configured forge host, got %s", url)
	}
}

func TestAuthService_Exchange(t *testing.T) {
	t.Run("exchanges a valid code for a token", func(t *testing.T) {
		forgeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"real-token","token_type":"bearer"}`))
		}))
		defer forgeServer.Close()

		svc := NewAuthService(svcmocks.NewMockUserRepository(t), nil, "client-id", "client-secret", "gitea", forgeServer.URL, "https://abyss.example.com/callback", "secret", "owner")

		token, err := svc.Exchange(context.Background(), "valid-code")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token.AccessToken != "real-token" {
			t.Fatalf("unexpected token: %+v", token)
		}
	})

	t.Run("propagates forge rejection of the code", func(t *testing.T) {
		forgeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
		}))
		defer forgeServer.Close()

		svc := NewAuthService(svcmocks.NewMockUserRepository(t), nil, "client-id", "client-secret", "gitea", forgeServer.URL, "https://abyss.example.com/callback", "secret", "owner")

		_, err := svc.Exchange(context.Background(), "bad-code")
		if err == nil {
			t.Fatal("expected an error for a rejected code")
		}
	})
}

func TestAuthService_GetForgeUser(t *testing.T) {
	t.Run("errors when forge is unavailable", func(t *testing.T) {
		svc := NewAuthService(svcmocks.NewMockUserRepository(t), nil, "", "", "", "", "", "secret", "owner")

		_, err := svc.GetForgeUser(context.Background(), &oauth2.Token{AccessToken: "x"})
		if err == nil {
			t.Fatal("expected an error when forge is nil")
		}
	})
}

func TestAuthService_GetOrCreateUser(t *testing.T) {
	t.Run("errors when forge is unavailable", func(t *testing.T) {
		svc := NewAuthService(svcmocks.NewMockUserRepository(t), nil, "", "", "", "", "", "secret", "owner")

		_, err := svc.GetOrCreateUser(context.Background(), &model.ForgeUser{ID: 1, Username: "alice"})
		if err == nil {
			t.Fatal("expected an error when forge is nil")
		}
	})

	t.Run("returns ErrUnauthorized when not a member of the owner org", func(t *testing.T) {
		forge := mocks.NewMockForge(t)
		forge.EXPECT().IsMemberOfOwner(context.Background(), "owner", "alice").Return(false, nil)

		svc := NewAuthService(svcmocks.NewMockUserRepository(t), forge, "", "", "", "", "", "secret", "owner")

		_, err := svc.GetOrCreateUser(context.Background(), &model.ForgeUser{ID: 1, Username: "alice"})
		if !errors.Is(err, ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("updates existing user's avatar and email", func(t *testing.T) {
		forge := mocks.NewMockForge(t)
		forge.EXPECT().IsMemberOfOwner(context.Background(), "owner", "alice").Return(true, nil)

		userRepo := svcmocks.NewMockUserRepository(t)
		existing := &model.User{ID: 9, ForgeID: 1, Username: "alice"}
		userRepo.EXPECT().GetByForgeID(context.Background(), int64(1)).Return(existing, nil)
		userRepo.EXPECT().Update(context.Background(), existing).Return(nil)

		svc := NewAuthService(userRepo, forge, "", "", "", "", "", "secret", "owner")

		got, err := svc.GetOrCreateUser(context.Background(), &model.ForgeUser{ID: 1, Username: "alice", Email: "alice@example.com", AvatarURL: "https://x/a.png"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Email != "alice@example.com" || got.AvatarURL != "https://x/a.png" {
			t.Fatalf("expected user fields refreshed, got %+v", got)
		}
	})

	t.Run("first ever user becomes admin", func(t *testing.T) {
		forge := mocks.NewMockForge(t)
		forge.EXPECT().IsMemberOfOwner(context.Background(), "owner", "alice").Return(true, nil)

		userRepo := svcmocks.NewMockUserRepository(t)
		userRepo.EXPECT().GetByForgeID(context.Background(), int64(1)).Return(nil, nil)
		userRepo.EXPECT().Count(context.Background()).Return(int64(0), nil)
		userRepo.On("Save", mock.Anything, mock.MatchedBy(func(u *model.User) bool { return u.IsAdmin })).Return(nil)

		svc := NewAuthService(userRepo, forge, "", "", "", "", "", "secret", "owner")

		got, err := svc.GetOrCreateUser(context.Background(), &model.ForgeUser{ID: 1, Username: "alice"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !got.IsAdmin {
			t.Fatalf("expected the first user to be admin, got %+v", got)
		}
	})

	t.Run("subsequent users are not admin", func(t *testing.T) {
		forge := mocks.NewMockForge(t)
		forge.EXPECT().IsMemberOfOwner(context.Background(), "owner", "bob").Return(true, nil)

		userRepo := svcmocks.NewMockUserRepository(t)
		userRepo.EXPECT().GetByForgeID(context.Background(), int64(2)).Return(nil, nil)
		userRepo.EXPECT().Count(context.Background()).Return(int64(1), nil)
		userRepo.On("Save", mock.Anything, mock.MatchedBy(func(u *model.User) bool { return !u.IsAdmin })).Return(nil)

		svc := NewAuthService(userRepo, forge, "", "", "", "", "", "secret", "owner")

		got, err := svc.GetOrCreateUser(context.Background(), &model.ForgeUser{ID: 2, Username: "bob"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.IsAdmin {
			t.Fatalf("expected a subsequent user to not be admin, got %+v", got)
		}
	})
}

func TestAuthService_JWTRoundtrip(t *testing.T) {
	svc := NewAuthService(svcmocks.NewMockUserRepository(t), nil, "", "", "", "", "", "session-secret", "owner")

	token, err := svc.GenerateJWT(&model.User{ID: 3, Username: "carol", IsAdmin: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims, err := svc.ValidateJWT(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.UserID != 3 || claims.Username != "carol" || !claims.IsAdmin {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestAuthService_PersonalTokenLifecycle(t *testing.T) {
	t.Run("generates and persists a hashed token, never the raw value", func(t *testing.T) {
		userRepo := svcmocks.NewMockUserRepository(t)
		user := &model.User{ID: 1}
		userRepo.EXPECT().Update(context.Background(), user).Return(nil)

		svc := NewAuthService(userRepo, nil, "", "", "", "", "", "secret", "owner")

		token, err := svc.GeneratePersonalToken(context.Background(), user)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token == "" {
			t.Fatal("expected a non-empty raw token")
		}
		if user.TokenHash == nil || *user.TokenHash == token {
			t.Fatalf("expected a stored hash distinct from the raw token, got %v", user.TokenHash)
		}
		if user.TokenLastEight != token[len(token)-8:] {
			t.Fatalf("expected TokenLastEight to match the token's suffix, got %q for token %q", user.TokenLastEight, token)
		}
	})

	t.Run("revokes a token", func(t *testing.T) {
		userRepo := svcmocks.NewMockUserRepository(t)
		hash := "existing-hash"
		user := &model.User{ID: 1, TokenHash: &hash, TokenLastEight: "abc12345"}
		userRepo.EXPECT().Update(context.Background(), user).Return(nil)

		svc := NewAuthService(userRepo, nil, "", "", "", "", "", "secret", "owner")

		if err := svc.RevokePersonalToken(context.Background(), user); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.TokenHash != nil || user.TokenLastEight != "" {
			t.Fatalf("expected token hash and suffix to be cleared, got hash=%v suffix=%q", user.TokenHash, user.TokenLastEight)
		}
	})
}

func TestAuthService_GetUserByToken(t *testing.T) {
	t.Run("finds the matching user among last-eight candidates", func(t *testing.T) {
		userRepo := svcmocks.NewMockUserRepository(t)
		user := &model.User{ID: 1}

		svc := NewAuthService(userRepo, nil, "", "", "", "", "", "secret", "owner")
		userRepo.EXPECT().Update(context.Background(), user).Return(nil)
		token, err := svc.GeneratePersonalToken(context.Background(), user)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		userRepo.EXPECT().GetByTokenLastEight(context.Background(), token[len(token)-8:]).Return([]model.User{*user}, nil)

		got, err := svc.GetUserByToken(context.Background(), token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil || got.ID != 1 {
			t.Fatalf("expected to find user by token, got %+v", got)
		}
	})

	t.Run("rejects a token whose hash does not match a candidate", func(t *testing.T) {
		userRepo := svcmocks.NewMockUserRepository(t)
		otherHash := "some-other-hash"
		candidate := model.User{ID: 2, TokenHash: &otherHash, TokenLastEight: "abc12345"}

		svc := NewAuthService(userRepo, nil, "", "", "", "", "", "secret", "owner")

		token := "0000000000000000000000000000000000000000000000000000000abc12345"
		userRepo.EXPECT().GetByTokenLastEight(context.Background(), "abc12345").Return([]model.User{candidate}, nil)

		got, err := svc.GetUserByToken(context.Background(), token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != nil {
			t.Fatalf("expected no match for an incorrect token, got %+v", got)
		}
	})

	t.Run("returns nil without a repository call for a too-short token", func(t *testing.T) {
		userRepo := svcmocks.NewMockUserRepository(t)
		svc := NewAuthService(userRepo, nil, "", "", "", "", "", "secret", "owner")

		got, err := svc.GetUserByToken(context.Background(), "short")
		if err != nil || got != nil {
			t.Fatalf("expected nil, nil for a short token, got %+v, %v", got, err)
		}
	})
}
