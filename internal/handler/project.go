package handler

import (
	"encoding/json"
	"net/http"

	"git.4thena.io/4thena/abys/internal/database"
	"git.4thena.io/4thena/abys/internal/dto/request"
	"git.4thena.io/4thena/abys/internal/dto/response"
	"git.4thena.io/4thena/abys/internal/model"
	"git.4thena.io/4thena/abys/internal/repository"
	"git.4thena.io/4thena/abys/internal/service"
	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	service service.ProjectService
}

func NewProjectHandler() *ProjectHandler {
	repository := repository.NewProjectRepository(database.Connection)
	service := service.NewProjectService(*repository)
	return &ProjectHandler{*service}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req request.CreateProject
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	project, err := h.service.SaveProject(r.Context(), &model.Project{
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
	projects, err := h.service.GetAllProjects(r.Context())
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

func (h *ProjectHandler) GetProjectByName(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	project, err := h.service.GetProjectByName(r.Context(), name)
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

func (h *ProjectHandler) RegisterProjectRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/projects").Subrouter()

	userRouter.HandleFunc("", h.CreateProject).Methods("POST")
	userRouter.HandleFunc("", h.GetAllProjects).Methods("GET")
	userRouter.HandleFunc("/{value}", h.GetProjectByName).Methods("GET")
}
