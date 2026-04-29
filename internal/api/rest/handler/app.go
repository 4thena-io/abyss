package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/4thena-io/abyss/internal/api/rest/request"
	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/go-chi/chi/v5"
)

type AppHandler struct {
	service           service.AppService
	deploymentService service.DeploymentService
}

func NewAppHandler(appService *service.AppService, deploymentService *service.DeploymentService) *AppHandler {
	return &AppHandler{*appService, *deploymentService}
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

func (h *AppHandler) GetAppDeployments(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	deployments, err := h.deploymentService.GetDeploymentsByApp(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := make([]response.Deployment, len(deployments))
	for i, d := range deployments {
		res[i] = response.Deployment{
			ID:          d.ID,
			AppID:       d.AppID,
			Environment: d.Environment,
			Status:      d.Status,
			Commit:      d.Commit,
			TriggeredBy: d.TriggeredBy,
			Duration:    d.Duration,
			DeployedAt:  d.DeployedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
