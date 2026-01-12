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

type AppHandler struct {
	service service.AppService
}

func NewAppHandler() *AppHandler {
	appRepository := repository.NewAppRepository(database.Connection)
	projectRepository := repository.NewProjectRepository(database.Connection)

	forge, err := forge.NewForge()
	if err != nil {
		log.Fatalf("failed to create a git provider: %s", err)
	}
	ciProvider, err := ci.NewCi()
	if err != nil {
		log.Fatalf("failed to create a ci provider: %s", err)
	}

	service := service.NewAppService(*appRepository, *projectRepository, forge, ciProvider)

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
	json.NewEncoder(w).Encode(response.App{
		ID:          app.ID,
		Name:        app.Name,
		Description: app.Description,
		Kind:        app.Kind,
		Language:    app.Language,
		RepoURL:     app.RepoURL,
		CiURL:       app.CiURL,
		ProjectID:   app.ProjectID,
		TemplateID:  app.TemplateID,
	})
}

func (h *AppHandler) GetAllApps(w http.ResponseWriter, r *http.Request) {
	apps, err := h.service.GetAllApps(r.Context())
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

func (h *AppHandler) GetAppByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	app, err := h.service.GetAppByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if app == nil {
		http.Error(w, "App not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.App{
		ID:          app.ID,
		Name:        app.Name,
		Description: app.Description,
		Kind:        app.Kind,
		Language:    app.Language,
		RepoURL:     app.RepoURL,
		CiURL:       app.CiURL,
		ProjectID:   app.ProjectID,
		TemplateID:  app.TemplateID,
	})
}

func (h *AppHandler) GetAppBuilds(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	builds, err := h.service.GetAppBuilds(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(builds)
}

func (h *AppHandler) RegisterRoutes(router *mux.Router) {
	appRouter := router.PathPrefix("/apps").Subrouter()

	appRouter.HandleFunc("", h.CreateApp).Methods("POST")
	appRouter.HandleFunc("", h.GetAllApps).Methods("GET")
	appRouter.HandleFunc("/{id}", h.GetAppByID).Methods("GET")
	appRouter.HandleFunc("/{id}/builds", h.GetAppBuilds).Methods("GET")
}
