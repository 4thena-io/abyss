package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/4thena-io/abyss/internal/api/rest/request"
	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/database"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/repository"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	projectService service.ProjectService
	appService     service.AppService
}

func NewProjectHandler() *ProjectHandler {
	projectRepository := repository.NewProjectRepository(database.Connection)
	appRepository := repository.NewAppRepository(database.Connection)
	templateRepository := repository.NewTemplateRepository(database.Connection)

	projectService := service.NewProjectService(*projectRepository)
	appService := service.NewAppServiceReadOnly(*appRepository, *projectRepository, *templateRepository)

	return &ProjectHandler{
		projectService: *projectService,
		appService:     *appService,
	}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req request.CreateProject
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	project, err := h.projectService.SaveProject(r.Context(), &model.Project{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.Project{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
	})
}

func (h *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.GetAllProjects(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := make([]response.Project, len(projects))
	for i, project := range projects {
		res[i] = response.Project{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	project, err := h.projectService.GetProjectByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if project == nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.Project{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
	})
}

func (h *ProjectHandler) GetProjectApps(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	apps, err := h.appService.GetAppsByProject(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.projectService.DeleteProject(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
