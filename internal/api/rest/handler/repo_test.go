package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
)

func TestRepoHandler_GetAllRepos(t *testing.T) {
	t.Run("returns repos as json", func(t *testing.T) {
		forge := &fakeForge{
			GetOrgReposFn: func(ctx context.Context, name string) ([]model.Repo, error) {
				return []model.Repo{{ID: 1, Name: "app", FullName: "owner/app", URL: "https://forge/owner/app"}}, nil
			},
		}
		svc := service.NewRepoService(forge, "owner")
		h := NewRepoHandler(svc)

		r := requestWithParams(http.MethodGet, "/api/repos", "", nil, nil)
		w := httptest.NewRecorder()
		h.GetAllRepos(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var got []response.Repo
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(got) != 1 || got[0].FullName != "owner/app" {
			t.Fatalf("unexpected repos: %+v", got)
		}
	})

	t.Run("returns 500 when the forge lookup fails", func(t *testing.T) {
		forge := &fakeForge{
			GetOrgReposFn: func(ctx context.Context, name string) ([]model.Repo, error) { return nil, errors.New("boom") },
		}
		svc := service.NewRepoService(forge, "owner")
		h := NewRepoHandler(svc)

		r := requestWithParams(http.MethodGet, "/api/repos", "", nil, nil)
		w := httptest.NewRecorder()
		h.GetAllRepos(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}
