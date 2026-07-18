package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/4thena-io/abyss/internal/config"
)

func TestSetupHandler_Status(t *testing.T) {
	t.Run("configured instance omits defaults", func(t *testing.T) {
		h := NewSetupHandler("", true, nil, make(chan struct{}, 1))

		r := requestWithParams(http.MethodGet, "/api/setup/status", "", nil, nil)
		w := httptest.NewRecorder()
		h.Status(w, r)

		var got map[string]any
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got["configured"] != true {
			t.Fatalf("expected configured=true, got %+v", got)
		}
		if _, ok := got["defaults"]; ok {
			t.Fatalf("expected no defaults for a configured instance, got %+v", got)
		}
	})

	t.Run("unconfigured instance echoes non-secret defaults", func(t *testing.T) {
		defaults := &config.Config{
			Database: config.DatabaseConfig{Type: "sqlite", Path: "./abyss.db"},
			Forge:    config.ForgeConfig{Type: "gitea", Host: "https://git.example.com", Owner: "acme"},
			CI:       config.CIConfig{Type: "woodpecker", Host: "https://ci.example.com"},
			Auth:     config.AuthConfig{ClientID: "client-123"},
		}
		h := NewSetupHandler("", false, defaults, make(chan struct{}, 1))

		r := requestWithParams(http.MethodGet, "/api/setup/status", "", nil, nil)
		w := httptest.NewRecorder()
		h.Status(w, r)

		var got map[string]any
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got["configured"] != false {
			t.Fatalf("expected configured=false, got %+v", got)
		}
		d, ok := got["defaults"].(map[string]any)
		if !ok {
			t.Fatalf("expected a defaults object, got %+v", got)
		}
		if d["forge_type"] != "gitea" || d["forge_owner"] != "acme" || d["client_id"] != "client-123" {
			t.Fatalf("unexpected defaults: %+v", d)
		}
	})
}

func TestSetupHandler_Configure(t *testing.T) {
	t.Run("returns 409 when already configured", func(t *testing.T) {
		h := NewSetupHandler("", true, nil, make(chan struct{}, 1))

		r := requestWithParams(http.MethodPost, "/api/setup", `{}`, nil, nil)
		w := httptest.NewRecorder()
		h.Configure(w, r)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})

	t.Run("returns 400 for malformed body", func(t *testing.T) {
		h := NewSetupHandler("", false, nil, make(chan struct{}, 1))

		r := requestWithParams(http.MethodPost, "/api/setup", "{not json", nil, nil)
		w := httptest.NewRecorder()
		h.Configure(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 when required fields are missing", func(t *testing.T) {
		h := NewSetupHandler("", false, nil, make(chan struct{}, 1))

		r := requestWithParams(http.MethodPost, "/api/setup", `{}`, nil, nil)
		w := httptest.NewRecorder()
		h.Configure(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 for postgres without host/user/name", func(t *testing.T) {
		h := NewSetupHandler("", false, nil, make(chan struct{}, 1))

		body := `{"db_type":"postgres","forge_type":"gitea","forge_host":"h","client_id":"c","client_secret":"s","ci_type":"woodpecker","ci_host":"h","ci_token":"t"}`
		r := requestWithParams(http.MethodPost, "/api/setup", body, nil, nil)
		w := httptest.NewRecorder()
		h.Configure(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("returns 400 for woodpecker without ci_host/ci_token", func(t *testing.T) {
		h := NewSetupHandler("", false, nil, make(chan struct{}, 1))

		body := `{"forge_type":"gitea","forge_host":"h","client_id":"c","client_secret":"s","ci_type":"woodpecker"}`
		r := requestWithParams(http.MethodPost, "/api/setup", body, nil, nil)
		w := httptest.NewRecorder()
		h.Configure(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("writes the config file, returns 201, and triggers a restart", func(t *testing.T) {
		configPath := filepath.Join(t.TempDir(), "abyss.yml")
		restartCh := make(chan struct{}, 1)
		h := NewSetupHandler(configPath, false, nil, restartCh)

		body := `{"forge_type":"gitea","forge_host":"https://git.example.com","client_id":"c","client_secret":"s","ci_type":"github-actions"}`
		r := requestWithParams(http.MethodPost, "/api/setup", body, nil, nil)
		w := httptest.NewRecorder()
		h.Configure(w, r)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("expected config file to be written: %v", err)
		}
		content := string(data)
		if !strings.Contains(content, `type: "gitea"`) || !strings.Contains(content, `host: "https://git.example.com"`) {
			t.Fatalf("unexpected config content: %s", content)
		}

		select {
		case <-restartCh:
		case <-time.After(time.Second):
			t.Fatal("expected a restart signal within 1s")
		}
	})
}

func TestBuildDatabaseYAML(t *testing.T) {
	t.Run("defaults to sqlite with the given path", func(t *testing.T) {
		got := buildDatabaseYAML(setupRequest{DBType: "sqlite", DBPath: "./abyss.db"})
		if !strings.Contains(got, `type: "sqlite"`) || !strings.Contains(got, `path: "./abyss.db"`) {
			t.Fatalf("unexpected yaml: %s", got)
		}
	})

	t.Run("defaults postgres port to 5432", func(t *testing.T) {
		got := buildDatabaseYAML(setupRequest{DBType: "postgres", DBHost: "h", DBUser: "u", DBName: "n"})
		if !strings.Contains(got, `port: "5432"`) {
			t.Fatalf("expected default postgres port 5432, got: %s", got)
		}
	})

	t.Run("defaults mysql port to 3306", func(t *testing.T) {
		got := buildDatabaseYAML(setupRequest{DBType: "mysql", DBHost: "h", DBUser: "u", DBName: "n"})
		if !strings.Contains(got, `port: "3306"`) {
			t.Fatalf("expected default mysql port 3306, got: %s", got)
		}
	})

	t.Run("respects an explicit port", func(t *testing.T) {
		got := buildDatabaseYAML(setupRequest{DBType: "postgres", DBHost: "h", DBUser: "u", DBName: "n", DBPort: "6543"})
		if !strings.Contains(got, `port: "6543"`) {
			t.Fatalf("expected explicit port 6543, got: %s", got)
		}
	})
}

func TestForgeBranch(t *testing.T) {
	if forgeBranch(setupRequest{}) != "main" {
		t.Fatal("expected default branch main")
	}
	if forgeBranch(setupRequest{ForgeBranch: "develop"}) != "develop" {
		t.Fatal("expected explicit branch to be respected")
	}
}

func TestYS_EscapesQuotes(t *testing.T) {
	got := ys(`has "quotes" inside`)
	want := `"has \"quotes\" inside"`
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
