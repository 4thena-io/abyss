package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type DocsHandler struct {
	docsService service.DocsService
}

func NewDocsHandler(docsService *service.DocsService) *DocsHandler {
	return &DocsHandler{*docsService}
}

func (h *DocsHandler) GetAppDocs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	docs, err := h.docsService.GetDocs(r.Context(), uint(id))
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, "docs not found")
		return
	}
	if err != nil {
		log.Error().Err(err).Uint64("app_id", id).Msg("failed to get docs")
		response.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}
