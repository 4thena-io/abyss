package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4thena-io/abyss/internal/service"
)

func TestHealthHandler_Check(t *testing.T) {
	forge := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer forge.Close()

	svc := service.NewHealthService(forge.URL, "github-actions", "")
	h := NewHealthHandler(svc)

	r := requestWithParams(http.MethodGet, "/api/health", "", nil, nil)
	w := httptest.NewRecorder()
	h.Check(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var status service.HealthStatus
	if err := json.NewDecoder(w.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !status.Forge.OK || !status.CI.OK {
		t.Fatalf("expected both healthy, got %+v", status)
	}
}
