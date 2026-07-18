package woodpecker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func writeJSON(t *testing.T, w http.ResponseWriter, status int, body string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func TestWoodpeckerCi_ActivateRepo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/repos" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("forge_remote_id"); got != "42" {
			t.Fatalf("expected forge_remote_id=42, got %q", got)
		}
		writeJSON(t, w, 200, `{"id":7,"full_name":"acme/my-app"}`)
	}))
	defer srv.Close()

	ci, err := NewWoodpeckerCi(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create ci client: %v", err)
	}

	got, err := ci.ActivateRepo(context.Background(), 42, "acme/my-app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 7 || got.Slug != "acme/my-app" || got.URL != srv.URL+"/repos/7" {
		t.Fatalf("unexpected ci repo: %+v", got)
	}
}

func TestWoodpeckerCi_DeleteRepo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/repos/7" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	ci, err := NewWoodpeckerCi(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create ci client: %v", err)
	}

	if err := ci.DeleteRepo(context.Background(), 7, "acme/my-app"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWoodpeckerCi_GetBuilds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/repos/7/pipelines" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		writeJSON(t, w, 200, `[
			{"id":1,"number":10,"status":"success","branch":"main","commit":"abc123"},
			{"id":2,"number":9,"status":"failure","branch":"main","commit":"def456"}
		]`)
	}))
	defer srv.Close()

	ci, err := NewWoodpeckerCi(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create ci client: %v", err)
	}

	got, err := ci.GetBuilds(context.Background(), 7, "acme/my-app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 builds, got %d", len(got))
	}
	if got[0].Number != 10 || got[0].Status != "success" || got[0].Commit != "abc123" {
		t.Fatalf("unexpected first build: %+v", got[0])
	}
	if got[0].Link != srv.URL+"/repos/7/pipeline/10" {
		t.Fatalf("unexpected build link: %s", got[0].Link)
	}
}

func TestWoodpeckerCi_ActivateRepo_PropagatesError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	ci, err := NewWoodpeckerCi(srv.URL, "token")
	if err != nil {
		t.Fatalf("failed to create ci client: %v", err)
	}

	if _, err := ci.ActivateRepo(context.Background(), 42, "acme/my-app"); err == nil {
		t.Fatal("expected an error when the CI server returns 500")
	}
}
