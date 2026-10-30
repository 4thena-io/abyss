package handler

import (
	"encoding/json"
	"net/http"

	"github.com/4thena-io/abyss/internal/service"
)

type HealthHandler struct {
	service *service.HealthService
}

func NewHealthHandler(svc *service.HealthService) *HealthHandler {
	return &HealthHandler{svc}
}

func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	status, ready := h.service.Ready(r.Context())
	w.Header().Set("Content-Type", "application/json")
	if !ready {
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	_ = json.NewEncoder(w).Encode(status)
}

// Check retains the existing diagnostic endpoint's 200 response; load balancers
// must use /api/ready, and process supervisors must use /api/live.
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	status := h.service.Check(r.Context())
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}
