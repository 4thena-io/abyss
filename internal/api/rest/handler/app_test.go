package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/4thena-io/abyss/internal/api/rest/middleware"
	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/auth"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/go-chi/chi/v5"
)

func newAppHandler(appRepo *fakeAppRepository, projectRepo *fakeProjectRepository, deploymentRepo *fakeDeploymentRepository) *AppHandler {
	appSvc := service.NewAppService(appRepo, projectRepo, &fakeTemplateRepository{}, &fakeForge{}, &fakeCI{}, nil, "owner", "", "secret", "main")
	deploymentSvc := service.NewDeploymentService(deploymentRepo)
	return NewAppHandler(appSvc, deploymentSvc)
}

// requestWithParams builds a request carrying chi URL params and, optionally,
// authenticated claims in the context, mirroring what the router/middleware
// would inject before the handler runs.
func requestWithParams(method, target string, body string, claims *auth.Claims, params map[string]string) *http.Request {
	var r *http.Request
	if body != "" {
		r = httptest.NewRequest(method, target, strings.NewReader(body))
	} else {
		r = httptest.NewRequest(method, target, nil)
	}

	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	if claims != nil {
		ctx = context.WithValue(ctx, middleware.ContextKeyUser, claims)
	}
	return r.WithContext(ctx)
}

func TestAppHandler_GetAppByID(t *testing.T) {
	t.Run("returns 404 when app missing", func(t *testing.T) {
		appRepo := &fakeAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return nil, nil },
		}
		h := newAppHandler(appRepo, &fakeProjectRepository{}, &fakeDeploymentRepository{})

		r := requestWithParams(http.MethodGet, "/apps/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetAppByID(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns 400 for invalid id", func(t *testing.T) {
		h := newAppHandler(&fakeAppRepository{}, &fakeProjectRepository{}, &fakeDeploymentRepository{})

		r := requestWithParams(http.MethodGet, "/apps/abc", "", nil, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		h.GetAppByID(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns app as json", func(t *testing.T) {
		appRepo := &fakeAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) {
				return &model.App{ID: id, Name: "my-app"}, nil
			},
		}
		h := newAppHandler(appRepo, &fakeProjectRepository{}, &fakeDeploymentRepository{})

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
		h := newAppHandler(&fakeAppRepository{}, &fakeProjectRepository{}, &fakeDeploymentRepository{})

		body := `{"name":"app"}`
		r := requestWithParams(http.MethodPost, "/apps", body, &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.CreateApp(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("returns 400 for malformed json body", func(t *testing.T) {
		h := newAppHandler(&fakeAppRepository{}, &fakeProjectRepository{}, &fakeDeploymentRepository{})

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
		appRepo := &fakeAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) {
				return &model.App{ID: 1, CreatorID: 42}, nil
			},
		}
		h := newAppHandler(appRepo, &fakeProjectRepository{}, &fakeDeploymentRepository{})

		body := `{"name":"renamed","description":"d"}`
		r := requestWithParams(http.MethodPut, "/apps/1", body, &auth.Claims{UserID: 999, IsAdmin: false}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.UpdateApp(w, r)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("returns 404 when app missing", func(t *testing.T) {
		appRepo := &fakeAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return nil, nil },
		}
		h := newAppHandler(appRepo, &fakeProjectRepository{}, &fakeDeploymentRepository{})

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
		appRepo := &fakeAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return app, nil },
			UpdateFn:  func(ctx context.Context, a *model.App) error { return nil },
		}
		h := newAppHandler(appRepo, &fakeProjectRepository{}, &fakeDeploymentRepository{})

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
		h := newAppHandler(&fakeAppRepository{}, &fakeProjectRepository{}, &fakeDeploymentRepository{})

		r := requestWithParams(http.MethodDelete, "/apps/abc", "", &auth.Claims{UserID: 1}, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		h.DeleteApp(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 404 when app missing", func(t *testing.T) {
		appRepo := &fakeAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return nil, nil },
		}
		h := newAppHandler(appRepo, &fakeProjectRepository{}, &fakeDeploymentRepository{})

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
		deploymentRepo := &fakeDeploymentRepository{
			GetByAppFn: func(ctx context.Context, appID uint) ([]model.Deployment, error) {
				return []model.Deployment{{ID: 1, AppID: appID, Environment: "prod"}}, nil
			},
		}
		h := newAppHandler(&fakeAppRepository{}, &fakeProjectRepository{}, deploymentRepo)

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
