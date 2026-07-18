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

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	status := h.service.Check(r.Context())
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}
