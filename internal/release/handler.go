package release

import (
	"encoding/json"
	"net/http"

	"git.d4ramirez.com/project-abyss/abys-api/internal/database"
	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler() *Handler {
	repository := NewRepository(database.Connection)
	service := NewService(*repository)
	return &Handler{*service}
}

func (h *Handler) SaveRelease(w http.ResponseWriter, r *http.Request) {
	var req CreateReleaseRequestDto
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

func (h *Handler) GetAllReleases(w http.ResponseWriter, r *http.Request) {
	version, err := h.service.GetAllReleases(r.Context(), mux.Vars(r)["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(version)
}

func (h *Handler) GetReleaseByTag(w http.ResponseWriter, r *http.Request) {
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

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	appRouter := router.PathPrefix("/releases").Subrouter()

	appRouter.HandleFunc("", handler.SaveRelease).Methods("POST")
	appRouter.HandleFunc("", handler.GetAllReleases).Methods("GET")
	appRouter.HandleFunc("/{tag}", handler.GetReleaseByTag).Methods("GET")
}
