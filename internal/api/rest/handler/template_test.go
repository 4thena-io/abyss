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

func newTemplateHandler(t *testing.T, templateRepo *mocks.MockTemplateRepository, forge *forgemocks.MockForge) *TemplateHandler {
	t.Helper()
	templateSvc := service.NewTemplateService(templateRepo, forge, "owner")
	appSvc := service.NewAppService(mocks.NewMockAppRepository(t), mocks.NewMockProjectRepository(t), templateRepo, forge, cimocks.NewMockCI(t), nil, "owner", "", "secret", "main")
	return NewTemplateHandler(templateSvc, appSvc)
}

func TestTemplateHandler_CreateTemplate(t *testing.T) {
	t.Run("returns 409 when a repo-based template name already exists", func(t *testing.T) {
		repo := mocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByName(mock.Anything, "taken").Return(&model.Template{ID: 1, Name: "taken"}, nil)
		h := newTemplateHandler(t, repo, forgemocks.NewMockForge(t))

		body := `{"source":"repo","name":"taken","repoUrl":"https://forge/owner/taken"}`
		r := requestWithParams(http.MethodPost, "/templates", body, &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.CreateTemplate(w, r)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("creates a repo-based template and returns 201", func(t *testing.T) {
		repo := mocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByName(mock.Anything, "new-template").Return(nil, nil)
		repo.On("Save", mock.Anything, mock.AnythingOfType("*model.Template")).Return(nil)
		h := newTemplateHandler(t, repo, forgemocks.NewMockForge(t))

		body := `{"source":"repo","name":"new-template","kind":"service","language":"go","repoUrl":"https://forge/owner/new-template"}`
		r := requestWithParams(http.MethodPost, "/templates", body, &auth.Claims{UserID: 7}, nil)
		w := httptest.NewRecorder()
		h.CreateTemplate(w, r)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
		}
		var got response.Template
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.Name != "new-template" || got.CreatorID != 7 {
			t.Fatalf("unexpected template: %+v", got)
		}
	})

	t.Run("returns 400 for malformed body", func(t *testing.T) {
		h := newTemplateHandler(t, mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t))

		r := requestWithParams(http.MethodPost, "/templates", "{not json", &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.CreateTemplate(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}

func TestTemplateHandler_UpdateTemplate(t *testing.T) {
	t.Run("returns 404 when template missing", func(t *testing.T) {
		repo := mocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		h := newTemplateHandler(t, repo, forgemocks.NewMockForge(t))

		body := `{"name":"n","description":"d","kind":"k","language":"l"}`
		r := requestWithParams(http.MethodPut, "/templates/1", body, &auth.Claims{UserID: 1}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.UpdateTemplate(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns 403 when caller is not creator or admin", func(t *testing.T) {
		repo := mocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.Template{ID: 1, CreatorID: 42}, nil)
		h := newTemplateHandler(t, repo, forgemocks.NewMockForge(t))

		body := `{"name":"n","description":"d","kind":"k","language":"l"}`
		r := requestWithParams(http.MethodPut, "/templates/1", body, &auth.Claims{UserID: 999}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.UpdateTemplate(w, r)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})
}

func TestTemplateHandler_GetTemplateByID(t *testing.T) {
	t.Run("returns 404 when template missing", func(t *testing.T) {
		repo := mocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		h := newTemplateHandler(t, repo, forgemocks.NewMockForge(t))

		r := requestWithParams(http.MethodGet, "/templates/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetTemplateByID(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns template as json", func(t *testing.T) {
		repo := mocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.Template{ID: 1, Name: "template"}, nil)
		h := newTemplateHandler(t, repo, forgemocks.NewMockForge(t))

		r := requestWithParams(http.MethodGet, "/templates/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetTemplateByID(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var got response.Template
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.Name != "template" {
			t.Fatalf("unexpected template: %+v", got)
		}
	})
}

func TestTemplateHandler_DeleteTemplate(t *testing.T) {
	t.Run("returns 404 when template missing", func(t *testing.T) {
		repo := mocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		h := newTemplateHandler(t, repo, forgemocks.NewMockForge(t))

		r := requestWithParams(http.MethodDelete, "/templates/1", "", &auth.Claims{UserID: 1}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.DeleteTemplate(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("deletes and returns 204", func(t *testing.T) {
		repo := mocks.NewMockTemplateRepository(t)
		template := &model.Template{ID: 1, CreatorID: 7}
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(template, nil)
		repo.EXPECT().Delete(mock.Anything, template).Return(nil)
		h := newTemplateHandler(t, repo, forgemocks.NewMockForge(t))

		r := requestWithParams(http.MethodDelete, "/templates/1", "", &auth.Claims{UserID: 7}, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.DeleteTemplate(w, r)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d body=%s", w.Code, w.Body.String())
		}
	})
}
