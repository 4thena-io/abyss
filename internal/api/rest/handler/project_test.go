package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/auth"
	cimocks "github.com/4thena-io/abyss/internal/integration/ci/mocks"
	forgemocks "github.com/4thena-io/abyss/internal/integration/forge/mocks"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/4thena-io/abyss/internal/service/mocks"
	"github.com/stretchr/testify/mock"
)

func newProjectHandler(t *testing.T, projectRepo *mocks.MockProjectRepository, appRepo *mocks.MockAppRepository) *ProjectHandler {
	t.Helper()
	projectSvc := service.NewProjectService(projectRepo)
	appSvc := service.NewAppService(appRepo, projectRepo, mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t), nil, "owner", "", "secret", "main")
	return NewProjectHandler(projectSvc, appSvc)
}

func TestProjectHandler_CreateProject(t *testing.T) {
	t.Run("returns 409 when name already exists", func(t *testing.T) {
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByName(mock.Anything, "taken").Return(&model.Project{ID: 1, Name: "taken"}, nil)
		h := newProjectHandler(t, projectRepo, mocks.NewMockAppRepository(t))

		body := `{"name":"taken","description":"d"}`
		r := requestWithParams(http.MethodPost, "/projects", body, &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.CreateProject(w, r)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("returns 400 for malformed body", func(t *testing.T) {
		h := newProjectHandler(t, mocks.NewMockProjectRepository(t), mocks.NewMockAppRepository(t))

		r := requestWithParams(http.MethodPost, "/projects", "{not json", &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.CreateProject(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("creates project and returns 201", func(t *testing.T) {
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByName(mock.Anything, "new-project").Return(nil, nil)
		projectRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.Project")).Return(nil)
		h := newProjectHandler(t, projectRepo, mocks.NewMockAppRepository(t))

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
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		h := newProjectHandler(t, projectRepo, mocks.NewMockAppRepository(t))

		r := requestWithParams(http.MethodGet, "/projects/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetProjectByID(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns project as json", func(t *testing.T) {
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.Project{ID: 1, Name: "proj"}, nil)
		h := newProjectHandler(t, projectRepo, mocks.NewMockAppRepository(t))

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
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.Project{ID: 1, CreatorID: 42}, nil)
		h := newProjectHandler(t, projectRepo, mocks.NewMockAppRepository(t))

		r := requestWithParams(http.MethodDelete, "/projects/1", "", &auth.Claims{UserID: 999}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.DeleteProject(w, r)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("deletes project and returns 204", func(t *testing.T) {
		project := &model.Project{ID: 1, CreatorID: 7}
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(project, nil)
		projectRepo.EXPECT().Delete(mock.Anything, project).Return(nil)
		h := newProjectHandler(t, projectRepo, mocks.NewMockAppRepository(t))

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
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		h := newProjectHandler(t, projectRepo, mocks.NewMockAppRepository(t))

		r := requestWithParams(http.MethodGet, "/projects/1/apps", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetProjectApps(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns apps for project", func(t *testing.T) {
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.Project{ID: 1}, nil)
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByProject(mock.Anything, uint(1)).Return([]model.App{{ID: 1, Name: "app-a", ProjectID: ptrUint(1)}}, nil)
		h := newProjectHandler(t, projectRepo, appRepo)

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

func ptrUint(v uint) *uint { return &v }
