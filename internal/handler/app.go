package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"git.d4ramirez.com/project-abyss/abys-api/internal/database"
	"git.d4ramirez.com/project-abyss/abys-api/internal/dto"
	"git.d4ramirez.com/project-abyss/abys-api/internal/ci"
	"git.d4ramirez.com/project-abyss/abys-api/internal/forge"
	"git.d4ramirez.com/project-abyss/abys-api/internal/repository"
	"git.d4ramirez.com/project-abyss/abys-api/internal/service"
	"github.com/gorilla/mux"
)

type AppHandler struct {
	service service.AppService
}

func NewAppHandler() *AppHandler {
	repository := repository.NewAppRepository(database.Connection)
	forge, err := forge.NewForge()
	if err != nil {
		log.Fatalf("failed to create a git provider: %s", err)
	}
	ciProvider, err := ci.NewCi()
	if err != nil {
		log.Fatalf("failed to create a ci provider: %s", err)
	}
	service := service.NewAppService(*repository, forge, ciProvider)
	return &AppHandler{*service}
}

func (h *AppHandler) CreateApp(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAppRequestDto
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

func (h *AppHandler) GetAllApps(w http.ResponseWriter, r *http.Request) {
	apps, err := h.service.GetAllApps(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

func (h *AppHandler) GetAppByName(w http.ResponseWriter, r *http.Request) {
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

func (h *AppHandler) RegisterAppRoutes(router *mux.Router) {
	appRouter := router.PathPrefix("/apps").Subrouter()

	appRouter.HandleFunc("", h.CreateApp).Methods("POST")
	appRouter.HandleFunc("", h.GetAllApps).Methods("GET")
	appRouter.HandleFunc("/{name}", h.GetAppByName).Methods("GET")

	appWithNameRouter := appRouter.PathPrefix("/{name}").Subrouter()
	NewReleaseHandler().RegisterReleaseRoutes(appWithNameRouter)
}
