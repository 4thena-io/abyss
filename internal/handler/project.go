package handler

import (
	"encoding/json"
	"net/http"

	"git.d4ramirez.com/project-abyss/abys-api/internal/database"
	"git.d4ramirez.com/project-abyss/abys-api/internal/dto"
	"git.d4ramirez.com/project-abyss/abys-api/internal/repository"
	"git.d4ramirez.com/project-abyss/abys-api/internal/service"
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
	var req dto.CreateProjectRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	app, err := h.service.SaveProject(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app)
}

func (h *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	apps, err := h.service.GetAllProjects(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
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
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) RegisterProjectRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/projects").Subrouter()

	userRouter.HandleFunc("", h.CreateProject).Methods("POST")
	userRouter.HandleFunc("", h.GetAllProjects).Methods("GET")
	userRouter.HandleFunc("/{value}", h.GetProjectByName).Methods("GET")
}
