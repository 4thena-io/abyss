package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestHealthRejectsHTTPFailures(t *testing.T) {
	for _, code := range []int{401, 403, 404, 500, 502, 503} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(code) }))
		svc := NewHealthService(server.URL, "woodpecker", server.URL)
		status := svc.Check(context.Background())
		server.Close()
		if status.Forge.OK || status.CI.OK {
			t.Fatalf("HTTP %d must not be healthy", code)
		}
	}
}
func TestReadinessRecoversWithDependencies(t *testing.T) {
	var dbOK, forgeOK atomic.Bool
	check := func(flag *atomic.Bool) func(context.Context) error {
		return func(context.Context) error {
			if !flag.Load() {
				return errors.New("secret connection credentials")
			}
			return nil
		}
	}
	svc := NewHealthService("https://forge.example", "github-actions", "").WithDependencies(check(&dbOK), check(&forgeOK))
	status, ready := svc.Ready(context.Background())
	if ready || status.Database == nil || status.Database.OK || status.Database.Error != "temporarily unavailable" {
		t.Fatalf("%+v ready=%v", status, ready)
	}
	dbOK.Store(true)
	if _, ready = svc.Ready(context.Background()); ready {
		t.Fatal("forge still down")
	}
	forgeOK.Store(true)
	if _, ready = svc.Ready(context.Background()); !ready {
		t.Fatal("did not recover")
	}
	dbOK.Store(false)
	if _, ready = svc.Ready(context.Background()); ready {
		t.Fatal("database outage missed")
	}
}
