package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"git.4thena.io/4thena/abys/internal/api/rest/request"
	"git.4thena.io/4thena/abys/internal/api/rest/response"
	"git.4thena.io/4thena/abys/internal/config"
	"git.4thena.io/4thena/abys/internal/database"
	"git.4thena.io/4thena/abys/internal/integration/ci"
	"git.4thena.io/4thena/abys/internal/integration/forge"
	"git.4thena.io/4thena/abys/internal/integration/git"
	"git.4thena.io/4thena/abys/internal/model"
	"git.4thena.io/4thena/abys/internal/repository"
	"git.4thena.io/4thena/abys/internal/service"
	"github.com/go-chi/chi/v5"
)

type AppHandler struct {
	service service.AppService
}

func NewAppHandler() *AppHandler {
	appRepository := repository.NewAppRepository(database.Connection)
	projectRepository := repository.NewProjectRepository(database.Connection)
	templateRepository := repository.NewTemplateRepository(database.Connection)

	forge, err := forge.NewForge()
	if err != nil {
		log.Fatalf("failed to create a forge provider: %s", err)
	}
	ciProvider, err := ci.NewCi()
	if err != nil {
		log.Fatalf("failed to create a ci provider: %s", err)
	}
	gitClient := git.New(config.Environment.ForgeToken)

	service := service.NewAppService(*appRepository, *projectRepository, *templateRepository, forge, ciProvider, gitClient)

	return &AppHandler{*service}
}

func (h *AppHandler) CreateApp(w http.ResponseWriter, r *http.Request) {
	var req request.CreateApp
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var app *model.App
	var err error

	if req.RepoID != 0 {
		// Create from existing repo
		app, err = h.service.CreateAppFromRepo(r.Context(), &model.App{
			Name:        req.Name,
			Description: req.Description,
			Kind:        req.Kind,
			Language:    req.Language,
			ProjectID:   req.ProjectID,
			RepoID:      req.RepoID,
		})
	} else if req.TemplateID != 0 {
		// Create from template
		app, err = h.service.CreateAppFromTemplate(r.Context(), &model.App{
			Name:        req.Name,
			Description: req.Description,
			Kind:        req.Kind,
			Language:    req.Language,
			ProjectID:   req.ProjectID,
			TemplateID:  &req.TemplateID,
		})
	} else {
		http.Error(w, "Must specify either templateId or repoId", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.App{
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

func (h *AppHandler) GetAppByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
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
	})
}

func (h *AppHandler) GetAppBuilds(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	builds, err := h.service.GetAppBuilds(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := make([]response.Build, len(builds))
	for i, build := range builds {
		res[i] = response.Build{
			ID:       build.ID,
			Number:   build.Number,
			Status:   build.Status,
			Branch:   build.Branch,
			Commit:   build.Commit,
			Duration: build.Duration,
			Link:     build.Link,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *AppHandler) DeleteApp(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteApp(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
