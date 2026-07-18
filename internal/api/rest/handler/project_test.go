package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/auth"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
)

func newProjectHandler(projectRepo *fakeProjectRepository, appRepo *fakeAppRepository) *ProjectHandler {
	projectSvc := service.NewProjectService(projectRepo)
	appSvc := service.NewAppService(appRepo, projectRepo, &fakeTemplateRepository{}, &fakeForge{}, &fakeCI{}, nil, "owner", "", "secret", "main")
	return NewProjectHandler(projectSvc, appSvc)
}

func TestProjectHandler_CreateProject(t *testing.T) {
	t.Run("returns 409 when name already exists", func(t *testing.T) {
		projectRepo := &fakeProjectRepository{
			GetByNameFn: func(ctx context.Context, name string) (*model.Project, error) {
				return &model.Project{ID: 1, Name: name}, nil
			},
		}
		h := newProjectHandler(projectRepo, &fakeAppRepository{})

		body := `{"name":"taken","description":"d"}`
		r := requestWithParams(http.MethodPost, "/projects", body, &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.CreateProject(w, r)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("returns 400 for malformed body", func(t *testing.T) {
		h := newProjectHandler(&fakeProjectRepository{}, &fakeAppRepository{})

		r := requestWithParams(http.MethodPost, "/projects", "{not json", &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.CreateProject(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("creates project and returns 201", func(t *testing.T) {
		projectRepo := &fakeProjectRepository{
			GetByNameFn: func(ctx context.Context, name string) (*model.Project, error) { return nil, nil },
			SaveFn:      func(ctx context.Context, project *model.Project) error { project.ID = 1; return nil },
		}
		h := newProjectHandler(projectRepo, &fakeAppRepository{})

		body := `{"name":"new-project","description":"d"}`
		r := requestWithParams(http.MethodPost, "/projects", body, &auth.Claims{UserID: 7}, nil)
		w := httptest.NewRecorder()
		h.CreateProject(w, r)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
		}
		var got response.Project
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.Name != "new-project" || got.CreatorID != 7 {
			t.Fatalf("unexpected project: %+v", got)
		}
	})
}

func TestProjectHandler_GetProjectByID(t *testing.T) {
	t.Run("returns 404 when project missing", func(t *testing.T) {
		projectRepo := &fakeProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return nil, nil },
		}
		h := newProjectHandler(projectRepo, &fakeAppRepository{})

		r := requestWithParams(http.MethodGet, "/projects/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetProjectByID(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns project as json", func(t *testing.T) {
		projectRepo := &fakeProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) {
				return &model.Project{ID: id, Name: "proj"}, nil
			},
		}
		h := newProjectHandler(projectRepo, &fakeAppRepository{})

		r := requestWithParams(http.MethodGet, "/projects/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetProjectByID(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var got response.Project
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.Name != "proj" {
			t.Fatalf("expected name proj, got %q", got.Name)
		}
	})
}

func TestProjectHandler_DeleteProject(t *testing.T) {
	t.Run("returns 403 when caller is not creator or admin", func(t *testing.T) {
		projectRepo := &fakeProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) {
				return &model.Project{ID: 1, CreatorID: 42}, nil
			},
		}
		h := newProjectHandler(projectRepo, &fakeAppRepository{})

		r := requestWithParams(http.MethodDelete, "/projects/1", "", &auth.Claims{UserID: 999}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.DeleteProject(w, r)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("deletes project and returns 204", func(t *testing.T) {
		projectRepo := &fakeProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) {
				return &model.Project{ID: 1, CreatorID: 7}, nil
			},
			DeleteFn: func(ctx context.Context, project *model.Project) error { return nil },
		}
		h := newProjectHandler(projectRepo, &fakeAppRepository{})

		r := requestWithParams(http.MethodDelete, "/projects/1", "", &auth.Claims{UserID: 7}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.DeleteProject(w, r)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d body=%s", w.Code, w.Body.String())
		}
	})
}

func TestProjectHandler_GetProjectApps(t *testing.T) {
	t.Run("returns 404 when project missing", func(t *testing.T) {
		projectRepo := &fakeProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return nil, nil },
		}
		h := newProjectHandler(projectRepo, &fakeAppRepository{})

		r := requestWithParams(http.MethodGet, "/projects/1/apps", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetProjectApps(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns apps for project", func(t *testing.T) {
		projectRepo := &fakeProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return &model.Project{ID: id}, nil },
		}
		appRepo := &fakeAppRepository{
			GetByProjectFn: func(ctx context.Context, id uint) ([]model.App, error) {
				return []model.App{{ID: 1, Name: "app-a", ProjectID: &id}}, nil
			},
		}
		h := newProjectHandler(projectRepo, appRepo)

		r := requestWithParams(http.MethodGet, "/projects/1/apps", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetProjectApps(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
		}
		var got []response.App
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(got) != 1 || got[0].Name != "app-a" {
			t.Fatalf("unexpected apps: %+v", got)
		}
	})
}
