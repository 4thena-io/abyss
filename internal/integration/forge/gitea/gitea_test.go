package gitea

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newFixtureServer wraps handler with the Gitea SDK's version handshake,
// which every client construction performs against /api/v1/version before
// anything else. Verified against the real SDK (code.gitea.io/sdk/gitea)
// with a request-logging probe rather than assumed.
func newFixtureServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/version" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"version":"1.20.0"}`))
			return
		}
		handler(w, r)
	}))
}

func writeJSON(t *testing.T, w http.ResponseWriter, status int, body string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func TestGiteaForge_GetAuthenticatedUser(t *testing.T) {
	srv := newFixtureServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/user" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		writeJSON(t, w, 200, `{"id":7,"login":"alice","full_name":"Alice Example","email":"alice@example.com","avatar_url":"https://forge/avatar.png"}`)
	})
	defer srv.Close()

	forge, err := NewGiteaForge(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create forge: %v", err)
	}

	got, err := forge.GetAuthenticatedUser(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 7 || got.Username != "alice" || got.FullName != "Alice Example" || got.Email != "alice@example.com" {
		t.Fatalf("unexpected user: %+v", got)
	}
}

func TestGiteaForge_GetAuthenticatedUser_FallsBackToUsernameWhenNoFullName(t *testing.T) {
	srv := newFixtureServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, 200, `{"id":7,"login":"alice","full_name":"","email":"","avatar_url":""}`)
	})
	defer srv.Close()

	forge, err := NewGiteaForge(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create forge: %v", err)
	}

	got, err := forge.GetAuthenticatedUser(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.FullName != "alice" {
		t.Fatalf("expected FullName to fall back to username, got %q", got.FullName)
	}
}

func TestGiteaForge_GetUserByToken(t *testing.T) {
	var gotAuth string
	srv := newFixtureServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		writeJSON(t, w, 200, `{"id":3,"login":"carol","full_name":"Carol Example","email":"carol@example.com"}`)
	})
	defer srv.Close()

	// GetUserByToken constructs a fresh client with the given per-user token
	// rather than reusing the forge's own bot token — verify the request
	// actually carries that token, not the one NewGiteaForge was built with.
	forge, err := NewGiteaForge(srv.URL, "bot-token")
	if err != nil {
		t.Fatalf("failed to create forge: %v", err)
	}

	got, err := forge.GetUserByToken(context.Background(), "user-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Username != "carol" {
		t.Fatalf("unexpected user: %+v", got)
	}
	if gotAuth == "" || gotAuth == "token bot-token" {
		t.Fatalf("expected the request to authenticate with the user token, got Authorization=%q", gotAuth)
	}
}

func TestGiteaForge_IsMemberOfOwner(t *testing.T) {
	t.Run("204 means the user is a member", func(t *testing.T) {
		srv := newFixtureServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/orgs/acme/members/alice" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		defer srv.Close()

		forge, err := NewGiteaForge(srv.URL, "token")
		if err != nil {
			t.Fatalf("failed to create forge: %v", err)
		}

		member, err := forge.IsMemberOfOwner(context.Background(), "acme", "alice")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !member {
			t.Fatal("expected member to be true")
		}
	})

	t.Run("404 means the user is not a member", func(t *testing.T) {
		srv := newFixtureServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer srv.Close()

		forge, err := NewGiteaForge(srv.URL, "token")
		if err != nil {
			t.Fatalf("failed to create forge: %v", err)
		}

		member, err := forge.IsMemberOfOwner(context.Background(), "acme", "bob")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if member {
			t.Fatal("expected member to be false")
		}
	})
}

func TestGiteaForge_CreateRepo(t *testing.T) {
	srv := newFixtureServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/orgs/acme/repos" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if body["name"] != "my-app" {
			t.Fatalf("unexpected repo name in request: %+v", body)
		}
		writeJSON(t, w, 201, `{"id":99,"name":"my-app","full_name":"acme/my-app","html_url":"https://forge/acme/my-app","clone_url":"https://forge/acme/my-app.git"}`)
	})
	defer srv.Close()

	forge, err := NewGiteaForge(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create forge: %v", err)
	}

	got, err := forge.CreateRepo(context.Background(), "acme", "my-app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 99 || got.FullName != "acme/my-app" || got.CloneURL != "https://forge/acme/my-app.git" {
		t.Fatalf("unexpected repo: %+v", got)
	}
}

func TestGiteaForge_GetOrgRepos(t *testing.T) {
	srv := newFixtureServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/orgs/acme/repos" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		writeJSON(t, w, 200, `[
			{"id":1,"name":"app-a","full_name":"acme/app-a","html_url":"https://forge/acme/app-a","clone_url":"https://forge/acme/app-a.git"},
			{"id":2,"name":"app-b","full_name":"acme/app-b","html_url":"https://forge/acme/app-b","clone_url":"https://forge/acme/app-b.git"}
		]`)
	})
	defer srv.Close()

	forge, err := NewGiteaForge(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create forge: %v", err)
	}

	got, err := forge.GetOrgRepos(context.Background(), "acme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].Name != "app-a" || got[1].Name != "app-b" {
		t.Fatalf("unexpected repos: %+v", got)
	}
}

func TestGiteaForge_GetRepo(t *testing.T) {
	srv := newFixtureServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/repositories/42" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		writeJSON(t, w, 200, `{"id":42,"name":"my-app","full_name":"acme/my-app","html_url":"https://forge/acme/my-app","clone_url":"https://forge/acme/my-app.git"}`)
	})
	defer srv.Close()

	forge, err := NewGiteaForge(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create forge: %v", err)
	}

	got, err := forge.GetRepo(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 42 || got.Name != "my-app" {
		t.Fatalf("unexpected repo: %+v", got)
	}
}

func TestGiteaForge_DeleteRepo(t *testing.T) {
	srv := newFixtureServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/repos/acme/my-app" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer srv.Close()

	forge, err := NewGiteaForge(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create forge: %v", err)
	}

	if err := forge.DeleteRepo(context.Background(), "acme", "my-app"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGiteaForge_CreateWebhook(t *testing.T) {
	srv := newFixtureServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/repos/acme/my-app/hooks" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		config, _ := body["config"].(map[string]any)
		if config["url"] != "https://abyss.example.com/hooks/1" {
			t.Fatalf("unexpected callback url in request: %+v", body)
		}
		if body["branch_filter"] != "main" {
			t.Fatalf("unexpected branch filter: %+v", body)
		}
		writeJSON(t, w, 201, `{"id":5}`)
	})
	defer srv.Close()

	forge, err := NewGiteaForge(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create forge: %v", err)
	}

	err = forge.CreateWebhook(context.Background(), "acme", "my-app", "https://abyss.example.com/hooks/1", "secret", "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGiteaForge_DeleteWebhook(t *testing.T) {
	var deletedID int64 = -1
	srv := newFixtureServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/repos/acme/my-app/hooks":
			writeJSON(t, w, 200, `[
				{"id":1,"config":{"url":"https://abyss.example.com/hooks/1?access_token=secret"}},
				{"id":2,"config":{"url":"https://other.example.com/unrelated"}}
			]`)
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/repos/acme/my-app/hooks/1":
			deletedID = 1
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodDelete:
			t.Fatalf("unexpected delete of unrelated hook: %s", r.URL.Path)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	defer srv.Close()

	forge, err := NewGiteaForge(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create forge: %v", err)
	}

	// DeleteWebhook matches on the callback URL with query params stripped,
	// so this must match hook 1 despite the differing ?access_token= value.
	err = forge.DeleteWebhook(context.Background(), "acme", "my-app", "https://abyss.example.com/hooks/1?access_token=different")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deletedID != 1 {
		t.Fatalf("expected hook 1 to be deleted, got deletedID=%d", deletedID)
	}
}
