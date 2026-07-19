package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/auth"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/4thena-io/abyss/internal/service/mocks"
	"github.com/stretchr/testify/mock"
)

func newTeamHandler(t *testing.T, teamRepo *mocks.MockTeamRepository, projectRepo *mocks.MockProjectRepository, appRepo *mocks.MockAppRepository, userLookup *mocks.MockUserLookup) *TeamHandler {
	t.Helper()
	svc := service.NewTeamService(teamRepo, projectRepo, appRepo, userLookup)
	return NewTeamHandler(svc)
}

func TestTeamHandler_CreateTeam(t *testing.T) {
	t.Run("returns 409 when team name already exists", func(t *testing.T) {
		teamRepo := mocks.NewMockTeamRepository(t)
		teamRepo.EXPECT().GetByName(mock.Anything, "taken").Return(&model.Team{ID: 1, Name: "taken"}, nil)
		h := newTeamHandler(t, teamRepo, mocks.NewMockProjectRepository(t), mocks.NewMockAppRepository(t), mocks.NewMockUserLookup(t))

		body := `{"name":"taken","description":"d"}`
		r := requestWithParams(http.MethodPost, "/teams", body, &auth.Claims{UserID: 1}, nil)
		w := httptest.NewRecorder()
		h.CreateTeam(w, r)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("creates team and returns 201", func(t *testing.T) {
		teamRepo := mocks.NewMockTeamRepository(t)
		teamRepo.EXPECT().GetByName(mock.Anything, "new-team").Return(nil, nil)
		teamRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.Team")).Return(nil)
		teamRepo.On("SaveMember", mock.Anything, mock.AnythingOfType("*model.TeamMember")).Return(nil)
		h := newTeamHandler(t, teamRepo, mocks.NewMockProjectRepository(t), mocks.NewMockAppRepository(t), mocks.NewMockUserLookup(t))

		body := `{"name":"new-team","description":"d"}`
		r := requestWithParams(http.MethodPost, "/teams", body, &auth.Claims{UserID: 7}, nil)
		w := httptest.NewRecorder()
		h.CreateTeam(w, r)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
		}
		var got response.Team
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.Name != "new-team" {
			t.Fatalf("unexpected team: %+v", got)
		}
	})
}

func TestTeamHandler_GetTeamByID(t *testing.T) {
	t.Run("returns 404 when team missing", func(t *testing.T) {
		teamRepo := mocks.NewMockTeamRepository(t)
		teamRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		h := newTeamHandler(t, teamRepo, mocks.NewMockProjectRepository(t), mocks.NewMockAppRepository(t), mocks.NewMockUserLookup(t))

		r := requestWithParams(http.MethodGet, "/teams/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetTeamByID(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("returns team stats as json", func(t *testing.T) {
		teamRepo := mocks.NewMockTeamRepository(t)
		teamRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.Team{ID: 1, Name: "team"}, nil)
		teamRepo.EXPECT().CountMembers(mock.Anything, uint(1)).Return(int64(2), nil)
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().CountByTeam(mock.Anything, uint(1)).Return(3, nil)
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().CountByTeam(mock.Anything, uint(1)).Return(4, nil)
		h := newTeamHandler(t, teamRepo, projectRepo, appRepo, mocks.NewMockUserLookup(t))

		r := requestWithParams(http.MethodGet, "/teams/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.GetTeamByID(w, r)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var got response.Team
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.MemberCount != 2 || got.ProjectCount != 3 || got.AppCount != 4 {
			t.Fatalf("unexpected team stats: %+v", got)
		}
	})
}

func TestTeamHandler_DeleteTeam(t *testing.T) {
	t.Run("returns 404 when team missing", func(t *testing.T) {
		teamRepo := mocks.NewMockTeamRepository(t)
		teamRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		h := newTeamHandler(t, teamRepo, mocks.NewMockProjectRepository(t), mocks.NewMockAppRepository(t), mocks.NewMockUserLookup(t))

		r := requestWithParams(http.MethodDelete, "/teams/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.DeleteTeam(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("deletes and returns 204", func(t *testing.T) {
		teamRepo := mocks.NewMockTeamRepository(t)
		team := &model.Team{ID: 1}
		teamRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(team, nil)
		teamRepo.EXPECT().Delete(mock.Anything, team).Return(nil)
		h := newTeamHandler(t, teamRepo, mocks.NewMockProjectRepository(t), mocks.NewMockAppRepository(t), mocks.NewMockUserLookup(t))

		r := requestWithParams(http.MethodDelete, "/teams/1", "", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.DeleteTeam(w, r)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})
}

func TestTeamHandler_AddTeamMember(t *testing.T) {
	t.Run("returns 404 when user does not exist", func(t *testing.T) {
		userLookup := mocks.NewMockUserLookup(t)
		userLookup.EXPECT().GetByUsername(mock.Anything, "ghost").Return(nil, nil)
		h := newTeamHandler(t, mocks.NewMockTeamRepository(t), mocks.NewMockProjectRepository(t), mocks.NewMockAppRepository(t), userLookup)

		body := `{"username":"ghost","role":"member"}`
		r := requestWithParams(http.MethodPost, "/teams/1/members", body, nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.AddTeamMember(w, r)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("defaults role to member and returns 201", func(t *testing.T) {
		teamRepo := mocks.NewMockTeamRepository(t)
		teamRepo.On("SaveMember", mock.Anything, mock.AnythingOfType("*model.TeamMember")).Return(nil)
		userLookup := mocks.NewMockUserLookup(t)
		userLookup.EXPECT().GetByUsername(mock.Anything, "alice").Return(&model.User{ID: 9, Username: "alice"}, nil)
		h := newTeamHandler(t, teamRepo, mocks.NewMockProjectRepository(t), mocks.NewMockAppRepository(t), userLookup)

		body := `{"username":"alice"}`
		r := requestWithParams(http.MethodPost, "/teams/1/members", body, nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.AddTeamMember(w, r)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
		}
		var got response.TeamMember
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.Role != "member" || got.Username != "alice" {
			t.Fatalf("unexpected member: %+v", got)
		}
	})
}
