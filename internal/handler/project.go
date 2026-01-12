package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"git.4thena.io/4thena/abys/internal/database"
	"git.4thena.io/4thena/abys/internal/dto/request"
	"git.4thena.io/4thena/abys/internal/dto/response"
	"git.4thena.io/4thena/abys/internal/integration/ci"
	"git.4thena.io/4thena/abys/internal/integration/forge"
	"git.4thena.io/4thena/abys/internal/model"
	"git.4thena.io/4thena/abys/internal/repository"
	"git.4thena.io/4thena/abys/internal/service"
	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	projectService service.ProjectService
	appService     service.AppService
}

func NewProjectHandler() *ProjectHandler {
	projectRepository := repository.NewProjectRepository(database.Connection)
	appRepository := repository.NewAppRepository(database.Connection)

	projectService := service.NewProjectService(*projectRepository)

	forge, err := forge.NewForge()
	if err != nil {
		log.Fatalf("failed to create a git provider: %s", err)
	}
	ciProvider, err := ci.NewCi()
	if err != nil {
		log.Fatalf("failed to create a ci provider: %s", err)
	}

	appService := service.NewAppService(*appRepository, *projectRepository, forge, ciProvider)

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
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
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
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
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
			ID:          app.ID,
			Name:        app.Name,
			Description: app.Description,
			Kind:        app.Kind,
			Language:    app.Language,
			RepoURL:     app.RepoURL,
			CiURL:       app.CiURL,
			ProjectID:   app.ProjectID,
			TemplateID:  app.TemplateID,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *ProjectHandler) RegisterRoutes(router *mux.Router) {
	projectRouter := router.PathPrefix("/projects").Subrouter()

	projectRouter.HandleFunc("", h.CreateProject).Methods("POST")
	projectRouter.HandleFunc("", h.GetAllProjects).Methods("GET")
	projectRouter.HandleFunc("/{id}", h.GetProjectByID).Methods("GET")
	projectRouter.HandleFunc("/{id}/apps", h.GetProjectApps).Methods("GET")
}
