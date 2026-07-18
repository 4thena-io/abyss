package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
)

func TestDocsHandler_GetAppDocs(t *testing.T) {
	t.Run("returns 400 for invalid id", func(t *testing.T) {
		svc := service.NewDocsService(&fakeAppRepository{}, nil, t.TempDir())
		h := NewDocsHandler(svc)

		r := requestWithParams(http.MethodGet, "/apps/abc/docs", "", nil, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		h.GetAppDocs(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 404 when app does not exist", func(t *testing.T) {
		appRepo := &fakeAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return nil, nil },
		}
		svc := service.NewDocsService(appRepo, nil, t.TempDir())
		h := NewDocsHandler(svc)

		r := requestWithParams(http.MethodGet, "/apps/1/docs", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetAppDocs(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns 404 when docs have not been rendered", func(t *testing.T) {
		appRepo := &fakeAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return &model.App{ID: id}, nil },
		}
		svc := service.NewDocsService(appRepo, nil, t.TempDir())
		h := NewDocsHandler(svc)

		r := requestWithParams(http.MethodGet, "/apps/1/docs", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetAppDocs(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns rendered docs as json", func(t *testing.T) {
		docsDir := t.TempDir()
		appRepo := &fakeAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return &model.App{ID: id}, nil },
		}
		appDocsDir := filepath.Join(docsDir, "1")
		if err := os.MkdirAll(appDocsDir, 0755); err != nil {
			t.Fatalf("failed to create docs dir: %v", err)
		}
		seed, err := json.Marshal(service.RenderedDocs{Pages: []service.DocPage{{Path: "index.md", HTML: "<h1>Hi</h1>"}}})
		if err != nil {
			t.Fatalf("failed to marshal seed docs: %v", err)
		}
		if err := os.WriteFile(filepath.Join(appDocsDir, "docs.json"), seed, 0644); err != nil {
			t.Fatalf("failed to write seed docs: %v", err)
		}

		svc := service.NewDocsService(appRepo, nil, docsDir)
		h := NewDocsHandler(svc)

		r := requestWithParams(http.MethodGet, "/apps/1/docs", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetAppDocs(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
		}
		var got service.RenderedDocs
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(got.Pages) != 1 || got.Pages[0].HTML != "<h1>Hi</h1>" {
			t.Fatalf("unexpected docs: %+v", got)
		}
	})
}
