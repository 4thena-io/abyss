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

func newAppHandler(t *testing.T, appRepo *mocks.MockAppRepository, projectRepo *mocks.MockProjectRepository, deploymentRepo *mocks.MockDeploymentRepository) *AppHandler {
	t.Helper()
	appSvc := service.NewAppService(appRepo, projectRepo, mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t), nil, "owner", "", "secret", "main")
	deploymentSvc := service.NewDeploymentService(deploymentRepo)
	return NewAppHandler(appSvc, deploymentSvc)
}

func TestAppHandler_GetAppByID(t *testing.T) {
	t.Run("returns 404 when app missing", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		h := newAppHandler(t, appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockDeploymentRepository(t))

		r := requestWithParams(http.MethodGet, "/apps/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetAppByID(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns 400 for invalid id", func(t *testing.T) {
		h := newAppHandler(t, mocks.NewMockAppRepository(t), mocks.NewMockProjectRepository(t), mocks.NewMockDeploymentRepository(t))

		r := requestWithParams(http.MethodGet, "/apps/abc", "", nil, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		h.GetAppByID(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns app as json", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.App{ID: 1, Name: "my-app"}, nil)
		h := newAppHandler(t, appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockDeploymentRepository(t))

		r := requestWithParams(http.MethodGet, "/apps/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetAppByID(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var got response.App
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.Name != "my-app" {
			t.Fatalf("expected name my-app, got %q", got.Name)
		}
	})
}

func TestAppHandler_CreateApp(t *testing.T) {
	t.Run("returns 400 when neither templateId nor repoId given", func(t *testing.T) {
		h := newAppHandler(t, mocks.NewMockAppRepository(t), mocks.NewMockProjectRepository(t), mocks.NewMockDeploymentRepository(t))

		body := `{"name":"app"}`
		r := requestWithParams(http.MethodPost, "/apps", body, &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.CreateApp(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("returns 400 for malformed json body", func(t *testing.T) {
		h := newAppHandler(t, mocks.NewMockAppRepository(t), mocks.NewMockProjectRepository(t), mocks.NewMockDeploymentRepository(t))

		r := requestWithParams(http.MethodPost, "/apps", "{not json", &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.CreateApp(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}

func TestAppHandler_UpdateApp(t *testing.T) {
	t.Run("returns 403 when caller is not creator or admin", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.App{ID: 1, CreatorID: 42}, nil)
		h := newAppHandler(t, appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockDeploymentRepository(t))

		body := `{"name":"renamed","description":"d"}`
		r := requestWithParams(http.MethodPut, "/apps/1", body, &auth.Claims{UserID: 999, IsAdmin: false}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.UpdateApp(w, r)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("returns 404 when app missing", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		h := newAppHandler(t, appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockDeploymentRepository(t))

		body := `{"name":"renamed","description":"d"}`
		r := requestWithParams(http.MethodPut, "/apps/1", body, &auth.Claims{UserID: 1}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.UpdateApp(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("updates app and returns 200", func(t *testing.T) {
		app := &model.App{ID: 1, CreatorID: 1, Name: "old"}
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(app, nil)
		appRepo.EXPECT().Update(mock.Anything, app).Return(nil)
		h := newAppHandler(t, appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockDeploymentRepository(t))

		body := `{"name":"renamed","description":"new desc"}`
		r := requestWithParams(http.MethodPut, "/apps/1", body, &auth.Claims{UserID: 1}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.UpdateApp(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
		}
		var got response.App
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.Name != "renamed" {
			t.Fatalf("expected name renamed, got %q", got.Name)
		}
	})
}

func TestAppHandler_DeleteApp(t *testing.T) {
	t.Run("returns 400 for invalid id", func(t *testing.T) {
		h := newAppHandler(t, mocks.NewMockAppRepository(t), mocks.NewMockProjectRepository(t), mocks.NewMockDeploymentRepository(t))

		r := requestWithParams(http.MethodDelete, "/apps/abc", "", &auth.Claims{UserID: 1}, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		h.DeleteApp(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 404 when app missing", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		h := newAppHandler(t, appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockDeploymentRepository(t))

		r := requestWithParams(http.MethodDelete, "/apps/1", "", &auth.Claims{UserID: 1}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.DeleteApp(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

func TestAppHandler_GetAppDeployments(t *testing.T) {
	t.Run("returns deployments as json", func(t *testing.T) {
		deploymentRepo := mocks.NewMockDeploymentRepository(t)
		deploymentRepo.EXPECT().GetByApp(mock.Anything, uint(1)).Return([]model.Deployment{{ID: 1, AppID: 1, Environment: "prod"}}, nil)
		h := newAppHandler(t, mocks.NewMockAppRepository(t), mocks.NewMockProjectRepository(t), deploymentRepo)

		r := requestWithParams(http.MethodGet, "/apps/1/deployments", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetAppDeployments(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
		}
		var got []response.Deployment
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(got) != 1 || got[0].Environment != "prod" {
			t.Fatalf("unexpected deployments: %+v", got)
		}
	})
}
