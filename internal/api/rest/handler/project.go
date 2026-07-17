package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/4thena-io/abyss/internal/api/rest/middleware"
	"github.com/4thena-io/abyss/internal/api/rest/request"
	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type ProjectHandler struct {
	projectService service.ProjectService
	appService     service.AppService
}

func NewProjectHandler(projectService *service.ProjectService, appService *service.AppService) *ProjectHandler {
	return &ProjectHandler{*projectService, *appService}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req request.CreateProject
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	defer r.Body.Close()

	claims := middleware.ClaimsFromContext(r.Context())
	project, err := h.projectService.SaveProject(r.Context(), &model.Project{
		Name:        req.Name,
		Description: req.Description,
		TeamID:      req.TeamID,
		CreatorID:   claims.UserID,
	})
	if errors.Is(err, service.ErrConflict) {
		response.Conflict(w, "project already exists")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to create project")
		response.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toProjectResponse(project))
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	var req request.UpdateProject
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	defer r.Body.Close()

	claims := middleware.ClaimsFromContext(r.Context())
	project, err := h.projectService.UpdateProject(r.Context(), uint(id), claims.UserID, claims.IsAdmin, req.Name, req.Description, req.TeamID)
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, "project not found")
		return
	}
	if errors.Is(err, service.ErrForbidden) {
		response.Forbidden(w)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to update project")
		response.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toProjectResponse(project))
}

func (h *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.GetAllProjects(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get projects")
		response.InternalError(w)
		return
	}

	res := make([]response.Project, len(projects))
	for i, project := range projects {
		res[i] = toProjectResponse(&project)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	project, err := h.projectService.GetProjectByID(r.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Msg("failed to get project")
		response.InternalError(w)
		return
	}
	if project == nil {
		response.NotFound(w, "project not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toProjectResponse(project))
}

func (h *ProjectHandler) GetProjectApps(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	apps, err := h.appService.GetAppsByProject(r.Context(), uint(id))
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, "project not found")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get project apps")
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

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	claims := middleware.ClaimsFromContext(r.Context())
	err = h.projectService.DeleteProject(r.Context(), uint(id), claims.UserID, claims.IsAdmin)
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, "project not found")
		return
	}
	if errors.Is(err, service.ErrForbidden) {
		response.Forbidden(w)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to delete project")
		response.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toProjectResponse(p *model.Project) response.Project {
	r := response.Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		TeamID:      p.TeamID,
		CreatorID:   p.CreatorID,
	}
	if p.Creator.Username != "" {
		r.CreatorUsername = p.Creator.Username
	}
	return r
}
