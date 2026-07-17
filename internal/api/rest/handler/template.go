package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/4thena-io/abyss/internal/api/rest/middleware"
	"github.com/4thena-io/abyss/internal/api/rest/request"
	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type TemplateHandler struct {
	service    service.TemplateService
	appService service.AppService
}

func NewTemplateHandler(templateService *service.TemplateService, appService *service.AppService) *TemplateHandler {
	return &TemplateHandler{*templateService, *appService}
}

func (h *TemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req request.CreateTemplate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	defer r.Body.Close()

	claims := middleware.ClaimsFromContext(r.Context())

	var template *model.Template
	var err error

	if req.Source == "blank" {
		template, err = h.service.CreateBlankTemplate(r.Context(), claims.UserID, req.Name, req.Description, req.Kind, req.Language)
	} else {
		template, err = h.service.SaveTemplate(r.Context(), &model.Template{
			Name:        req.Name,
			Description: req.Description,
			Kind:        req.Kind,
			Language:    req.Language,
			CloneURL:    req.RepoURL + ".git",
			RepoURL:     req.RepoURL,
			CreatorID:   claims.UserID,
		})
	}

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
	json.NewEncoder(w).Encode(toTemplateResponse(template))
}

func (h *TemplateHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	var req request.UpdateTemplate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	defer r.Body.Close()

	claims := middleware.ClaimsFromContext(r.Context())
	t, err := h.service.UpdateTemplate(r.Context(), uint(id), claims.UserID, claims.IsAdmin, req.Name, req.Description, req.Kind, req.Language)
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, "template not found")
		return
	}
	if errors.Is(err, service.ErrForbidden) {
		response.Forbidden(w)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to update template")
		response.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toTemplateResponse(t))
}

func (h *TemplateHandler) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.service.GetAllTemplates(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get templates")
		response.InternalError(w)
		return
	}

	res := make([]response.Template, len(templates))
	for i, t := range templates {
		res[i] = toTemplateResponse(&t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *TemplateHandler) GetTemplateByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	template, err := h.service.GetTemplateByID(r.Context(), uint(id))
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
	json.NewEncoder(w).Encode(toTemplateResponse(template))
}

func (h *TemplateHandler) GetTemplateApps(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	apps, err := h.appService.GetAppsByTemplate(r.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Msg("failed to get template apps")
		response.InternalError(w)
		return
	}

	res := make([]response.App, len(apps))
	for i, app := range apps {
		res[i] = toAppResponse(&app)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *TemplateHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	claims := middleware.ClaimsFromContext(r.Context())
	err = h.service.DeleteTemplate(r.Context(), uint(id), claims.UserID, claims.IsAdmin)
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, "template not found")
		return
	}
	if errors.Is(err, service.ErrForbidden) {
		response.Forbidden(w)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to delete template")
		response.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toTemplateResponse(t *model.Template) response.Template {
	r := response.Template{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Kind:        t.Kind,
		Language:    t.Language,
		RepoURL:     t.RepoURL,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		CreatorID:   t.CreatorID,
	}
	if t.Creator.Username != "" {
		r.CreatorUsername = t.Creator.Username
	}
	return r
}
