package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/4thena-io/abyss/internal/service/mocks"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_GetUser(t *testing.T) {
	t.Run("returns 400 for invalid id", func(t *testing.T) {
		teamSvc := service.NewTeamService(mocks.NewMockTeamRepository(t), mocks.NewMockProjectRepository(t), mocks.NewMockAppRepository(t), mocks.NewMockUserLookup(t))
		h := NewUserHandler(teamSvc)

		r := requestWithParams(http.MethodGet, "/users/abc", "", nil, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		h.GetUser(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("returns 404 when user does not exist", func(t *testing.T) {
		userLookup := mocks.NewMockUserLookup(t)
		userLookup.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		teamSvc := service.NewTeamService(mocks.NewMockTeamRepository(t), mocks.NewMockProjectRepository(t), mocks.NewMockAppRepository(t), userLookup)
		h := NewUserHandler(teamSvc)

		r := requestWithParams(http.MethodGet, "/users/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetUser(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns user profile with aggregated team stats", func(t *testing.T) {
		userLookup := mocks.NewMockUserLookup(t)
		userLookup.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.User{ID: 1, Username: "alice", IsAdmin: true}, nil)

		teamRepo := mocks.NewMockTeamRepository(t)
		teamRepo.EXPECT().GetByUserID(mock.Anything, uint(1)).Return([]model.TeamMember{{TeamID: 5, UserID: 1, Role: "owner"}}, nil)
		teamRepo.EXPECT().GetByID(mock.Anything, uint(5)).Return(&model.Team{ID: 5, Name: "platform"}, nil)
		teamRepo.EXPECT().CountMembers(mock.Anything, uint(5)).Return(int64(3), nil)

		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().CountByTeam(mock.Anything, uint(5)).Return(2, nil)
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().CountByTeam(mock.Anything, uint(5)).Return(4, nil)

		teamSvc := service.NewTeamService(teamRepo, projectRepo, appRepo, userLookup)
		h := NewUserHandler(teamSvc)

		r := requestWithParams(http.MethodGet, "/users/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetUser(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
		}
		var got response.UserProfile
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.Username != "alice" || !got.IsAdmin || got.TeamCount != 1 {
			t.Fatalf("unexpected profile: %+v", got)
		}
		if len(got.Teams) != 1 || got.Teams[0].Name != "platform" || got.Teams[0].Role != "owner" {
			t.Fatalf("unexpected teams: %+v", got.Teams)
		}
	})
}
