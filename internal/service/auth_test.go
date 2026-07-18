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
	t.Run("generates and persists a token", func(t *testing.T) {
		userRepo := svcmocks.NewMockUserRepository(t)
		user := &model.User{ID: 1}
		userRepo.EXPECT().Update(context.Background(), user).Return(nil)

		svc := NewAuthService(userRepo, nil, "", "", "", "", "", "secret", "owner")

		token, err := svc.GeneratePersonalToken(context.Background(), user)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token == "" || user.Token == nil || *user.Token != token {
			t.Fatalf("expected token to be generated and attached to the user, got token=%q user.Token=%v", token, user.Token)
		}
	})

	t.Run("revokes a token", func(t *testing.T) {
		userRepo := svcmocks.NewMockUserRepository(t)
		tok := "existing-token"
		user := &model.User{ID: 1, Token: &tok}
		userRepo.EXPECT().Update(context.Background(), user).Return(nil)

		svc := NewAuthService(userRepo, nil, "", "", "", "", "", "secret", "owner")

		if err := svc.RevokePersonalToken(context.Background(), user); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Token != nil {
			t.Fatalf("expected token to be cleared, got %v", user.Token)
		}
	})
}
