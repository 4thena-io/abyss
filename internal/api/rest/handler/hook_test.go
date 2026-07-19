package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/4thena-io/abyss/internal/service"
	"github.com/4thena-io/abyss/internal/service/mocks"
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

func giteaSignedRequest(secret, body, target string, params map[string]string) *http.Request {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	sig := hex.EncodeToString(mac.Sum(nil))

	r := requestWithParams(http.MethodPost, target, body, nil, params)
	r.Header.Set("X-Gitea-Signature", sig)
	return r
}

func TestHookHandler_Forge(t *testing.T) {
	t.Run("returns 401 when no secret is configured (fail closed)", func(t *testing.T) {
		docsSvc := service.NewDocsService(mocks.NewMockAppRepository(t), nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "", "gitea")

		r := requestWithParams(http.MethodPost, "/api/hooks/forge/1", `{"ref":"refs/heads/main"}`, nil, map[string]string{"appID": "1"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("returns 401 when the signature does not match the secret", func(t *testing.T) {
		docsSvc := service.NewDocsService(mocks.NewMockAppRepository(t), nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "expected-secret", "gitea")

		body := `{"ref":"refs/heads/main"}`
		r := giteaSignedRequest("wrong-secret", body, "/api/hooks/forge/1", map[string]string{"appID": "1"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("returns 401 when the signature header is missing", func(t *testing.T) {
		docsSvc := service.NewDocsService(mocks.NewMockAppRepository(t), nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "expected-secret", "gitea")

		r := requestWithParams(http.MethodPost, "/api/hooks/forge/1", `{"ref":"refs/heads/main"}`, nil, map[string]string{"appID": "1"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("accepts a GitLab-style token header", func(t *testing.T) {
		docsSvc := service.NewDocsService(mocks.NewMockAppRepository(t), nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "expected-secret", "gitlab")

		r := requestWithParams(http.MethodPost, "/api/hooks/forge/1", `{"ref":"refs/heads/develop"}`, nil, map[string]string{"appID": "1"})
		r.Header.Set("X-Gitlab-Token", "expected-secret")
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 (skipped, wrong branch), got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("returns 400 for invalid app id", func(t *testing.T) {
		docsSvc := service.NewDocsService(mocks.NewMockAppRepository(t), nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "secret", "gitea")

		r := giteaSignedRequest("secret", `{"ref":"refs/heads/main"}`, "/api/hooks/forge/abc", map[string]string{"appID": "abc"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 400 for malformed payload", func(t *testing.T) {
		docsSvc := service.NewDocsService(mocks.NewMockAppRepository(t), nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "secret", "gitea")

		r := giteaSignedRequest("secret", "{not json", "/api/hooks/forge/1", map[string]string{"appID": "1"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("skips with 200 when push is on a different branch", func(t *testing.T) {
		docsSvc := service.NewDocsService(mocks.NewMockAppRepository(t), nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "secret", "gitea")

		r := giteaSignedRequest("secret", `{"ref":"refs/heads/develop"}`, "/api/hooks/forge/1", map[string]string{"appID": "1"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 (skipped), got %d", w.Code)
		}
	})

	t.Run("skips with 200 when no docs-relevant files changed", func(t *testing.T) {
		docsSvc := service.NewDocsService(mocks.NewMockAppRepository(t), nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "secret", "gitea")

		body := `{"ref":"refs/heads/main","commits":[{"modified":["src/main.go"]}]}`
		r := giteaSignedRequest("secret", body, "/api/hooks/forge/1", map[string]string{"appID": "1"})
		w := httptest.NewRecorder()
		h.Forge(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 (skipped), got %d", w.Code)
		}
	})

	t.Run("returns 500 when a docs-relevant push fails to render (git not configured)", func(t *testing.T) {
		// RenderDocs rejects with a nil git client before ever touching the
		// app repository, so no fake app data is needed here.
		docsSvc := service.NewDocsService(mocks.NewMockAppRepository(t), nil, t.TempDir())
		h := NewHookHandler(docsSvc, "main", "secret", "gitea")

		body := `{"ref":"refs/heads/main","commits":[{"modified":["docs/index.md"]}]}`
		r := giteaSignedRequest("secret", body, "/api/hooks/forge/1", map[string]string{"appID": "1"})
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
