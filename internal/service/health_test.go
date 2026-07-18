package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthService_Check(t *testing.T) {
	t.Run("forge and self-hosted ci reachable", func(t *testing.T) {
		forge := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
		defer forge.Close()
		ci := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
		defer ci.Close()

		svc := NewHealthService(forge.URL, "woodpecker", ci.URL)
		status := svc.Check(context.Background())

		if !status.Forge.OK || !status.CI.OK {
			t.Fatalf("expected both healthy, got %+v", status)
		}
	})

	t.Run("unreachable host reports not ok", func(t *testing.T) {
		svc := NewHealthService("http://127.0.0.1:1", "woodpecker", "http://127.0.0.1:1")
		status := svc.Check(context.Background())

		if status.Forge.OK || status.Forge.Error == "" {
			t.Fatalf("expected forge to be unreachable, got %+v", status.Forge)
		}
		if status.CI.OK || status.CI.Error == "" {
			t.Fatalf("expected ci to be unreachable, got %+v", status.CI)
		}
	})

	t.Run("unconfigured host reports not configured", func(t *testing.T) {
		svc := NewHealthService("", "woodpecker", "")
		status := svc.Check(context.Background())

		if status.Forge.OK || status.Forge.Error != "not configured" {
			t.Fatalf("expected 'not configured', got %+v", status.Forge)
		}
	})

	t.Run("hosted CI types skip the ping and report healthy", func(t *testing.T) {
		for _, ciType := range []string{"github-actions", "gitlab-ci"} {
			svc := NewHealthService("", ciType, "")
			status := svc.Check(context.Background())
			if !status.CI.OK {
				t.Fatalf("expected hosted CI type %q to report OK without a host, got %+v", ciType, status.CI)
			}
		}
	})

	t.Run("self-hosted CI without a configured host reports not configured", func(t *testing.T) {
		svc := NewHealthService("", "woodpecker", "")
		status := svc.Check(context.Background())
		if status.CI.OK || status.CI.Error != "not configured" {
			t.Fatalf("expected 'not configured', got %+v", status.CI)
		}
	})
}
