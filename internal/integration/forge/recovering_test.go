package forge

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/4thena-io/abyss/internal/config"
)

func TestRecoveringForge(t *testing.T) {
	var online atomic.Bool
	var botAllowed atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !online.Load() {
			w.WriteHeader(502)
			return
		}
		switch r.URL.Path {
		case "/api/v1/version":
			_, _ = w.Write([]byte(`{"version":"10.0.0"}`))
		case "/api/v1/user":
			if r.Header.Get("Authorization") == "token bot" && !botAllowed.Load() {
				w.WriteHeader(401)
				return
			}
			_, _ = w.Write([]byte(`{"id":1,"login":"alice","full_name":"Alice","email":"alice@example.com"}`))
		default:
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	for _, kind := range []string{"gitea", "forgejo"} {
		t.Run(kind, func(t *testing.T) {
			online.Store(false)
			botAllowed.Store(false)
			f, err := NewRecovering(config.ForgeConfig{Type: kind, Host: server.URL, Token: "bot"})
			if err != nil {
				t.Fatal(err)
			}
			if _, err = f.GetAuthenticatedUser(context.Background()); !errors.Is(err, ErrUnavailable) {
				t.Fatalf("expected unavailable, got %v", err)
			}
			online.Store(true)
			user, err := f.GetUserByToken(context.Background(), "oauth")
			if err != nil || user.Username != "alice" {
				t.Fatalf("OAuth identity must not depend on bot: user=%v err=%v", user, err)
			}
			if _, err = f.GetAuthenticatedUser(context.Background()); err == nil {
				t.Fatal("bot should still fail")
			}
			botAllowed.Store(true)
			if _, err = f.GetAuthenticatedUser(context.Background()); err != nil {
				t.Fatalf("bot did not recover: %v", err)
			}
		})
	}
}
func TestRecoveringOAuthEndpointsOffline(t *testing.T) {
	f, err := NewRecovering(config.ForgeConfig{Type: "gitlab", Host: "http://127.0.0.1:1"})
	if err != nil {
		t.Fatal(err)
	}
	a, b := f.OAuthEndpoints()
	if a != "/oauth/authorize" || b != "/oauth/token" {
		t.Fatalf("%s %s", a, b)
	}
}
