package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/4thena-io/abyss/internal/service"
)

func TestPushPayload_IsOnBranch(t *testing.T) {
	p := pushPayload{Ref: "refs/heads/main"}
	if !p.isOnBranch("main") {
		t.Fatal("expected ref refs/heads/main to match branch main")
	}
	if p.isOnBranch("develop") {
		t.Fatal("expected ref refs/heads/main not to match branch develop")
	}
}

func TestPushPayload_TouchesDocs(t *testing.T) {
	t.Run("assumes docs changed when no file lists are present", func(t *testing.T) {
		p := pushPayload{Ref: "refs/heads/main"}
		if !p.touchesDocs() {
			t.Fatal("expected touchesDocs to default true when file info is missing")
		}
	})

	t.Run("true when a commit touches docs/", func(t *testing.T) {
		p := pushPayload{Commits: []struct {
			Added    []string `json:"added"`
			Modified []string `json:"modified"`
			Removed  []string `json:"removed"`
		}{{Added: []string{"docs/index.md"}}}}
		if !p.touchesDocs() {
			t.Fatal("expected touchesDocs to be true for a docs/ file")
		}
	})

	t.Run("true when a commit touches .abyss.yml", func(t *testing.T) {
		p := pushPayload{Commits: []struct {
			Added    []string `json:"added"`
			Modified []string `json:"modified"`
			Removed  []string `json:"removed"`
		}{{Modified: []string{".abyss.yml"}}}}
		if !p.touchesDocs() {
			t.Fatal("expected touchesDocs to be true for .abyss.yml")
		}
	})

	t.Run("false when file lists are present but unrelated", func(t *testing.T) {
		p := pushPayload{Commits: []struct {
			Added    []string `json:"added"`
			Modified []string `json:"modified"`
			Removed  []string `json:"removed"`
		}{{Modified: []string{"src/main.go"}}}}
		if p.touchesDocs() {
			t.Fatal("expected touchesDocs to be false when nothing relevant changed")
		}
	})
}

func TestHookHandler_Forge(t *testing.T) {
	t.Run("returns 401 when access_token does not match the webhook secret", func(t *testing.T) {
		docsSvc := service.NewDocsService(&fakeAppRepository{}, nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "expected-secret")

		r := requestWithParams(http.MethodPost, "/api/hooks/forge/1?access_token=wrong", `{"ref":"refs/heads/main"}`, nil, map[string]string{"appID": "1"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("returns 400 for invalid app id", func(t *testing.T) {
		docsSvc := service.NewDocsService(&fakeAppRepository{}, nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "")

		r := requestWithParams(http.MethodPost, "/api/hooks/forge/abc", `{"ref":"refs/heads/main"}`, nil, map[string]string{"appID": "abc"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 for malformed payload", func(t *testing.T) {
		docsSvc := service.NewDocsService(&fakeAppRepository{}, nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "")

		r := requestWithParams(http.MethodPost, "/api/hooks/forge/1", "{not json", nil, map[string]string{"appID": "1"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("skips with 200 when push is on a different branch", func(t *testing.T) {
		docsSvc := service.NewDocsService(&fakeAppRepository{}, nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "")

		r := requestWithParams(http.MethodPost, "/api/hooks/forge/1", `{"ref":"refs/heads/develop"}`, nil, map[string]string{"appID": "1"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 (skipped), got %d", w.Code)
		}
	})

	t.Run("skips with 200 when no docs-relevant files changed", func(t *testing.T) {
		docsSvc := service.NewDocsService(&fakeAppRepository{}, nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "")

		body := `{"ref":"refs/heads/main","commits":[{"modified":["src/main.go"]}]}`
		r := requestWithParams(http.MethodPost, "/api/hooks/forge/1", body, nil, map[string]string{"appID": "1"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 (skipped), got %d", w.Code)
		}
	})

	t.Run("returns 500 when a docs-relevant push fails to render (git not configured)", func(t *testing.T) {
		// RenderDocs rejects with a nil git client before ever touching the
		// app repository, so no fake app data is needed here.
		docsSvc := service.NewDocsService(&fakeAppRepository{}, nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "")

		body := `{"ref":"refs/heads/main","commits":[{"modified":["docs/index.md"]}]}`
		r := requestWithParams(http.MethodPost, "/api/hooks/forge/1", body, nil, map[string]string{"appID": "1"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500 since RenderDocs has no git client configured, got %d body=%s", w.Code, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), "internal error") {
			t.Fatalf("unexpected body: %s", w.Body.String())
		}
	})
}
