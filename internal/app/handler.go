package app

import (
	"encoding/json"
	"log"
	"net/http"

	"git.d4ramirez.com/project-abyss/abys-api/internal/database"
	"git.d4ramirez.com/project-abyss/abys-api/internal/integration/git"
	"git.d4ramirez.com/project-abyss/abys-api/internal/release"
	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler() *Handler {
	repository := NewRepository(database.Connection)
	gitProvider, err := git.NewGitProvider("gitea")
	if err != nil {
		log.Fatalf("failed to create a git provider: %s", err)
	}
	service := NewService(*repository, gitProvider)
	return &Handler{*service}
}

func (h *Handler) CreateApp(w http.ResponseWriter, r *http.Request) {
	var req CreateAppRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	app, err := h.service.CreateApp(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app)
}

func (h *Handler) GetAllApps(w http.ResponseWriter, r *http.Request) {
	apps, err := h.service.GetAllApps(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

func (h *Handler) GetAppByName(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	app, err := h.service.GetAppByName(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if app == nil {
		http.Error(w, "App not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	appRouter := router.PathPrefix("/apps").Subrouter()

	appRouter.HandleFunc("", h.CreateApp).Methods("POST")
	appRouter.HandleFunc("", h.GetAllApps).Methods("GET")
	appRouter.HandleFunc("/{name}", h.GetAppByName).Methods("GET")

	appWithNameRouter := appRouter.PathPrefix("/{name}").Subrouter()
	release.NewHandler().RegisterRoutes(appWithNameRouter)
}
