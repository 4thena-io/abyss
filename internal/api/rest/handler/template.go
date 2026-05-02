package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/4thena-io/abyss/internal/api/rest/request"
	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type TemplateHandler struct {
	service service.TemplateService
}

func NewTemplateHandler(service *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{*service}
}

func (h *TemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req request.CreateTemplate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	defer r.Body.Close()

	template, err := h.service.SaveTemplate(r.Context(), &model.Template{
		Name:        req.Name,
		Description: req.Description,
		Kind:        req.Kind,
		Language:    req.Language,
		CloneURL:    req.RepoURL + ".git",
		RepoURL:     req.RepoURL,
	})
	if errors.Is(err, service.ErrConflict) {
		response.Conflict(w, "template already exists")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to create template")
		response.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.Template{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		Kind:        template.Kind,
		Language:    template.Language,
		RepoURL:     template.RepoURL,
	})
}

func (h *TemplateHandler) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.service.GetAllTemplates(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get templates")
		response.InternalError(w)
		return
	}

	res := make([]response.Template, len(templates))
	for i, template := range templates {
		res[i] = response.Template{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			Kind:        template.Kind,
			Language:    template.Language,
			RepoURL:     template.RepoURL,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *TemplateHandler) GetTemplateByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "id")

	template, err := h.service.GetTemplateByName(r.Context(), name)
	if err != nil {
		log.Error().Err(err).Msg("failed to get template")
		response.InternalError(w)
		return
	}
	if template == nil {
		response.NotFound(w, "template not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.Template{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		Kind:        template.Kind,
		Language:    template.Language,
		RepoURL:     template.RepoURL,
	})
}

func (h *TemplateHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	err = h.service.DeleteTemplate(r.Context(), uint(id))
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, "template not found")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to delete template")
		response.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
