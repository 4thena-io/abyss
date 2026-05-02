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

	project, err := h.projectService.SaveProject(r.Context(), &model.Project{
		Name:        req.Name,
		Description: req.Description,
		TeamID:      req.TeamID,
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
	json.NewEncoder(w).Encode(response.Project{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		TeamID:      project.TeamID,
	})
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
		res[i] = response.Project{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			TeamID:      project.TeamID,
		}
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
	json.NewEncoder(w).Encode(response.Project{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		TeamID:      project.TeamID,
	})
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
		res[i] = response.App{
			ID:           app.ID,
			Name:         app.Name,
			Description:  app.Description,
			Kind:         app.Kind,
			Language:     app.Language,
			RepoFullName: app.RepoFullName,
			RepoURL:      app.RepoURL,
			CiURL:        app.CIURL,
			ProjectID:    app.ProjectID,
			TemplateID:   app.TemplateID,
		}
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

	err = h.projectService.DeleteProject(r.Context(), uint(id))
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, "project not found")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to delete project")
		response.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
