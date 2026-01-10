package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"git.4thena.io/4thena/abys/internal/ci"
	"git.4thena.io/4thena/abys/internal/database"
	"git.4thena.io/4thena/abys/internal/dto/request"
	"git.4thena.io/4thena/abys/internal/dto/response"
	"git.4thena.io/4thena/abys/internal/forge"
	"git.4thena.io/4thena/abys/internal/model"
	"git.4thena.io/4thena/abys/internal/repository"
	"git.4thena.io/4thena/abys/internal/service"
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
	var req request.CreateApp
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	appModel := &model.App{
		Name:        req.Name,
		Description: req.Description,
		Kind:        req.Kind,
		Language:    req.Language,
		ProjectID:   req.ProjectID,
		TemplateID:  req.TemplateID,
	}

	app, err := h.service.CreateApp(r.Context(), appModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.CreateApp{
		ID: app.ID,
		Name: app.Name,
		Description: app.Description,
		Kind: app.Kind,
		Language: app.Language,
		RepoURL: app.RepoURL,
		CiURL: app.CiURL,
	})
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
}
