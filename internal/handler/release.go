package handler

import (
	"encoding/json"
	"net/http"

	"git.4thena.io/4thena/abys/internal/database"
	"git.4thena.io/4thena/abys/internal/dto"
	"git.4thena.io/4thena/abys/internal/repository"
	"git.4thena.io/4thena/abys/internal/service"
	"github.com/gorilla/mux"
)

type ReleaseHandler struct {
	service service.ReleaseService
}

func NewReleaseHandler() *ReleaseHandler {
	repository := repository.NewReleaseRepository(database.Connection)
	service := service.NewReleaseService(*repository)
	return &ReleaseHandler{*service}
}

func (h *ReleaseHandler) SaveRelease(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateReleaseRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	version, err := h.service.SaveRelease(r.Context(), req, mux.Vars(r)["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(version)
}

func (h *ReleaseHandler) GetAllReleases(w http.ResponseWriter, r *http.Request) {
	version, err := h.service.GetAllReleases(r.Context(), mux.Vars(r)["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(version)
}

func (h *ReleaseHandler) GetReleaseByTag(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	tag := mux.Vars(r)["tag"]

	project, err := h.service.GetReleaseByTag(r.Context(), name, tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if project == nil {
		http.Error(w, "version not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (h *ReleaseHandler) RegisterReleaseRoutes(router *mux.Router) {
	appRouter := router.PathPrefix("/releases").Subrouter()

	appRouter.HandleFunc("", h.SaveRelease).Methods("POST")
	appRouter.HandleFunc("", h.GetAllReleases).Methods("GET")
	appRouter.HandleFunc("/{tag}", h.GetReleaseByTag).Methods("GET")
}
