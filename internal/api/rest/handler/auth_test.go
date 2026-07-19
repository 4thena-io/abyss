package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/4thena-io/abyss/internal/auth"
	forgemocks "github.com/4thena-io/abyss/internal/integration/forge/mocks"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
	svcmocks "github.com/4thena-io/abyss/internal/service/mocks"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Config(t *testing.T) {
	svc := service.NewAuthService(svcmocks.NewMockUserRepository(t), forgemocks.NewMockForge(t), "", "", "gitea", "https://forge.example.com", "", "secret", "owner")
	h := NewAuthHandler(svc)

	r := requestWithParams(http.MethodGet, "/api/auth/config", "", nil, nil)
	w := httptest.NewRecorder()
	h.Config(w, r)

	var got map[string]string
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got["forge_type"] != "gitea" || got["forge_host"] != "https://forge.example.com" {
		t.Fatalf("unexpected config: %+v", got)
	}
}

func TestAuthHandler_Login(t *testing.T) {
	svc := service.NewAuthService(svcmocks.NewMockUserRepository(t), forgemocks.NewMockForge(t), "client-id", "", "gitea", "https://forge.example.com", "https://abyss.example.com/callback", "secret", "owner")
	h := NewAuthHandler(svc)

	r := requestWithParams(http.MethodGet, "/api/auth/login", "", nil, nil)
	w := httptest.NewRecorder()
	h.Login(w, r)

	if w.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", w.Code)
	}
	location := w.Header().Get("Location")
	if !strings.HasPrefix(location, "https://forge.example.com/login/oauth/authorize") {
		t.Fatalf("unexpected redirect location: %s", location)
	}
	if len(w.Result().Cookies()) == 0 {
		t.Fatal("expected an oauth_state cookie to be set")
	}
}

func TestAuthHandler_Callback(t *testing.T) {
	t.Run("returns 400 when state does not match the cookie", func(t *testing.T) {
		svc := service.NewAuthService(svcmocks.NewMockUserRepository(t), forgemocks.NewMockForge(t), "", "", "", "", "", "secret", "owner")
		h := NewAuthHandler(svc)

		r := httptest.NewRequest(http.MethodGet, "/api/auth/callback?state=wrong", nil)
		r.AddCookie(&http.Cookie{Name: "oauth_state", Value: "expected"})
		w := httptest.NewRecorder()
		h.Callback(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("completes the flow and sets a session cookie on success", func(t *testing.T) {
		forgeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"real-token","token_type":"bearer"}`))
		}))
		defer forgeServer.Close()

		forge := forgemocks.NewMockForge(t)
		forge.EXPECT().GetUserByToken(mock.Anything, "real-token").Return(&model.ForgeUser{ID: 1, Username: "alice"}, nil)
		forge.EXPECT().IsMemberOfOwner(mock.Anything, "owner", "alice").Return(true, nil)

		userRepo := svcmocks.NewMockUserRepository(t)
		userRepo.EXPECT().GetByForgeID(mock.Anything, int64(1)).Return(nil, nil)
		userRepo.EXPECT().Count(mock.Anything).Return(int64(0), nil)
		userRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

		svc := service.NewAuthService(userRepo, forge, "client-id", "client-secret", "gitea", forgeServer.URL, "https://abyss.example.com/callback", "secret", "owner")
		h := NewAuthHandler(svc)

		r := httptest.NewRequest(http.MethodGet, "/api/auth/callback?state=abc123&code=valid-code", nil)
		r.AddCookie(&http.Cookie{Name: "oauth_state", Value: "abc123"})
		w := httptest.NewRecorder()
		h.Callback(w, r)

		if w.Code != http.StatusFound {
			t.Fatalf("expected 302, got %d body=%s", w.Code, w.Body.String())
		}
		var sessionCookie *http.Cookie
		for _, c := range w.Result().Cookies() {
			if c.Name == "session" {
				sessionCookie = c
			}
		}
		if sessionCookie == nil || sessionCookie.Value == "" {
			t.Fatal("expected a session cookie to be set")
		}
	})
}

func TestAuthHandler_Me(t *testing.T) {
	t.Run("returns 401 when unauthenticated", func(t *testing.T) {
		svc := service.NewAuthService(svcmocks.NewMockUserRepository(t), forgemocks.NewMockForge(t), "", "", "", "", "", "secret", "owner")
		h := NewAuthHandler(svc)

		r := requestWithParams(http.MethodGet, "/api/auth/me", "", nil, nil)
		w := httptest.NewRecorder()
		h.Me(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("returns the authenticated user's profile", func(t *testing.T) {
		userRepo := svcmocks.NewMockUserRepository(t)
		userRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.User{ID: 1, Username: "alice", Email: "a@example.com"}, nil)
		svc := service.NewAuthService(userRepo, forgemocks.NewMockForge(t), "", "", "gitea", "https://forge.example.com", "", "secret", "owner")
		h := NewAuthHandler(svc)

		r := requestWithParams(http.MethodGet, "/api/auth/me", "", &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.Me(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var got map[string]any
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got["username"] != "alice" {
			t.Fatalf("unexpected profile: %+v", got)
		}
	})
}

func TestAuthHandler_TokenStatus(t *testing.T) {
	t.Run("reports has_token false when no token is set", func(t *testing.T) {
		userRepo := svcmocks.NewMockUserRepository(t)
		userRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.User{ID: 1}, nil)
		svc := service.NewAuthService(userRepo, forgemocks.NewMockForge(t), "", "", "", "", "", "secret", "owner")
		h := NewAuthHandler(svc)

		r := requestWithParams(http.MethodGet, "/api/auth/token", "", &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.TokenStatus(w, r)

		var got map[string]bool
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got["has_token"] {
			t.Fatalf("expected has_token false, got %+v", got)
		}
	})
}

func TestAuthHandler_GenerateToken(t *testing.T) {
	userRepo := svcmocks.NewMockUserRepository(t)
	user := &model.User{ID: 1}
	userRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(user, nil)
	userRepo.On("Update", mock.Anything, user).Return(nil)
	svc := service.NewAuthService(userRepo, forgemocks.NewMockForge(t), "", "", "", "", "", "secret", "owner")
	h := NewAuthHandler(svc)

	r := requestWithParams(http.MethodPost, "/api/auth/token", "", &auth.Claims{UserID: 1}, nil)
	w := httptest.NewRecorder()
	h.GenerateToken(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
	var got map[string]string
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got["token"] == "" {
		t.Fatalf("expected a non-empty token, got %+v", got)
	}
}

func TestAuthHandler_RevokeToken(t *testing.T) {
	userRepo := svcmocks.NewMockUserRepository(t)
	tok := "existing"
	user := &model.User{ID: 1, Token: &tok}
	userRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(user, nil)
	userRepo.On("Update", mock.Anything, user).Return(nil)
	svc := service.NewAuthService(userRepo, forgemocks.NewMockForge(t), "", "", "", "", "", "secret", "owner")
	h := NewAuthHandler(svc)

	r := requestWithParams(http.MethodDelete, "/api/auth/token", "", &auth.Claims{UserID: 1}, nil)
	w := httptest.NewRecorder()
	h.RevokeToken(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	svc := service.NewAuthService(svcmocks.NewMockUserRepository(t), forgemocks.NewMockForge(t), "", "", "", "", "", "secret", "owner")
	h := NewAuthHandler(svc)

	r := requestWithParams(http.MethodPost, "/api/auth/logout", "", nil, nil)
	w := httptest.NewRecorder()
	h.Logout(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	var sessionCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "session" {
			sessionCookie = c
		}
	}
	if sessionCookie == nil || sessionCookie.MaxAge >= 0 {
		t.Fatalf("expected the session cookie to be cleared, got %+v", sessionCookie)
	}
}
