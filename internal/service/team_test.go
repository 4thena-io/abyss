package service

import (
	"context"
	"errors"
	"testing"

	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service/mocks"
)

func newTestTeamService(t *testing.T) (*TeamService, *mocks.MockTeamRepository, *mocks.MockProjectRepository, *mocks.MockAppRepository, *mocks.MockUserLookup) {
	t.Helper()
	teamRepo := mocks.NewMockTeamRepository(t)
	projectRepo := mocks.NewMockProjectRepository(t)
	appRepo := mocks.NewMockAppRepository(t)
	userRepo := mocks.NewMockUserLookup(t)
	return NewTeamService(teamRepo, projectRepo, appRepo, userRepo), teamRepo, projectRepo, appRepo, userRepo
}

func TestTeamService_SaveTeam(t *testing.T) {
	t.Run("returns ErrConflict when name already exists", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		teamRepo.EXPECT().GetByName(context.Background(), "taken").Return(&model.Team{ID: 1, Name: "taken"}, nil)

		_, err := svc.SaveTeam(context.Background(), &model.Team{Name: "taken"}, 1)
		if !errors.Is(err, ErrConflict) {
			t.Fatalf("expected ErrConflict, got %v", err)
		}
	})

	t.Run("creates team and adds creator as owner", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		team := &model.Team{Name: "new-team"}
		teamRepo.EXPECT().GetByName(context.Background(), "new-team").Return(nil, nil)
		teamRepo.EXPECT().Save(context.Background(), team).Return(nil)
		teamRepo.EXPECT().SaveMember(context.Background(), &model.TeamMember{TeamID: team.ID, UserID: 7, Role: "owner"}).Return(nil)

		got, err := svc.SaveTeam(context.Background(), team, 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != team {
			t.Fatalf("expected the same team instance to be returned")
		}
	})

	t.Run("propagates member save failure", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		team := &model.Team{Name: "new-team"}
		teamRepo.EXPECT().GetByName(context.Background(), "new-team").Return(nil, nil)
		teamRepo.EXPECT().Save(context.Background(), team).Return(nil)
		teamRepo.EXPECT().SaveMember(context.Background(), &model.TeamMember{TeamID: team.ID, UserID: 7, Role: "owner"}).Return(errBoom)

		_, err := svc.SaveTeam(context.Background(), team, 7)
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})
}

func TestTeamService_UpdateTeam(t *testing.T) {
	t.Run("returns ErrNotFound when team missing", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		teamRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(nil, nil)

		_, err := svc.UpdateTeam(context.Background(), 1, 1, false, "new", "desc")
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden when caller is not a member", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		team := &model.Team{ID: 1, Name: "old"}
		teamRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(team, nil)
		teamRepo.EXPECT().GetMemberByUserID(context.Background(), uint(1), uint(999)).Return(nil, nil)

		_, err := svc.UpdateTeam(context.Background(), 1, 999, false, "renamed", "desc")
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("returns ErrForbidden when caller is a member but not an owner", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		team := &model.Team{ID: 1, Name: "old"}
		teamRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(team, nil)
		teamRepo.EXPECT().GetMemberByUserID(context.Background(), uint(1), uint(999)).Return(&model.TeamMember{TeamID: 1, UserID: 999, Role: "member"}, nil)

		_, err := svc.UpdateTeam(context.Background(), 1, 999, false, "renamed", "desc")
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("allows a platform admin who is not a member", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		team := &model.Team{ID: 1, Name: "old"}
		teamRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(team, nil)
		teamRepo.EXPECT().Update(context.Background(), team).Return(nil)

		got, err := svc.UpdateTeam(context.Background(), 1, 999, true, "renamed", "new desc")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "renamed" || got.Description != "new desc" {
			t.Fatalf("team not updated as expected: %+v", got)
		}
	})

	t.Run("allows a team owner", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		team := &model.Team{ID: 1, Name: "old"}
		teamRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(team, nil)
		teamRepo.EXPECT().GetMemberByUserID(context.Background(), uint(1), uint(7)).Return(&model.TeamMember{TeamID: 1, UserID: 7, Role: "owner"}, nil)
		teamRepo.EXPECT().Update(context.Background(), team).Return(nil)

		got, err := svc.UpdateTeam(context.Background(), 1, 7, false, "renamed", "new desc")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "renamed" || got.Description != "new desc" {
			t.Fatalf("team not updated as expected: %+v", got)
		}
	})
}

func TestTeamService_DeleteTeam(t *testing.T) {
	t.Run("returns ErrNotFound when team missing", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		teamRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(nil, nil)

		err := svc.DeleteTeam(context.Background(), 1, 1, false)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden for a non-owner member", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		team := &model.Team{ID: 1}
		teamRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(team, nil)
		teamRepo.EXPECT().GetMemberByUserID(context.Background(), uint(1), uint(999)).Return(&model.TeamMember{TeamID: 1, UserID: 999, Role: "member"}, nil)

		err := svc.DeleteTeam(context.Background(), 1, 999, false)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("deletes existing team for the owner", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		team := &model.Team{ID: 1}
		teamRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(team, nil)
		teamRepo.EXPECT().GetMemberByUserID(context.Background(), uint(1), uint(7)).Return(&model.TeamMember{TeamID: 1, UserID: 7, Role: "owner"}, nil)
		teamRepo.EXPECT().Delete(context.Background(), team).Return(nil)

		if err := svc.DeleteTeam(context.Background(), 1, 7, false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestTeamService_AddTeamMember(t *testing.T) {
	t.Run("returns ErrForbidden when caller is not an owner", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		teamRepo.EXPECT().GetMemberByUserID(context.Background(), uint(1), uint(999)).Return(&model.TeamMember{TeamID: 1, UserID: 999, Role: "member"}, nil)

		_, err := svc.AddTeamMember(context.Background(), 1, 999, false, "alice", "member")
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("returns ErrNotFound when user does not exist", func(t *testing.T) {
		svc, teamRepo, _, _, userRepo := newTestTeamService(t)
		teamRepo.EXPECT().GetMemberByUserID(context.Background(), uint(1), uint(7)).Return(&model.TeamMember{TeamID: 1, UserID: 7, Role: "owner"}, nil)
		userRepo.EXPECT().GetByUsername(context.Background(), "ghost").Return(nil, nil)

		_, err := svc.AddTeamMember(context.Background(), 1, 7, false, "ghost", "member")
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("owner adds member and attaches user", func(t *testing.T) {
		svc, teamRepo, _, _, userRepo := newTestTeamService(t)
		teamRepo.EXPECT().GetMemberByUserID(context.Background(), uint(1), uint(7)).Return(&model.TeamMember{TeamID: 1, UserID: 7, Role: "owner"}, nil)
		user := &model.User{ID: 9, Username: "alice"}
		userRepo.EXPECT().GetByUsername(context.Background(), "alice").Return(user, nil)
		teamRepo.EXPECT().SaveMember(context.Background(), &model.TeamMember{TeamID: 1, UserID: 9, Role: "member"}).Return(nil)

		got, err := svc.AddTeamMember(context.Background(), 1, 7, false, "alice", "member")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.User.Username != "alice" {
			t.Fatalf("expected attached user, got %+v", got.User)
		}
	})

	t.Run("platform admin can add a member without being one", func(t *testing.T) {
		svc, teamRepo, _, _, userRepo := newTestTeamService(t)
		user := &model.User{ID: 9, Username: "alice"}
		userRepo.EXPECT().GetByUsername(context.Background(), "alice").Return(user, nil)
		teamRepo.EXPECT().SaveMember(context.Background(), &model.TeamMember{TeamID: 1, UserID: 9, Role: "member"}).Return(nil)

		got, err := svc.AddTeamMember(context.Background(), 1, 999, true, "alice", "member")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.User.Username != "alice" {
			t.Fatalf("expected attached user, got %+v", got.User)
		}
	})
}

func TestTeamService_RemoveTeamMember(t *testing.T) {
	t.Run("returns ErrForbidden when caller is not an owner", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		teamRepo.EXPECT().GetMemberByUserID(context.Background(), uint(1), uint(999)).Return(nil, nil)

		err := svc.RemoveTeamMember(context.Background(), 1, 5, 999, false)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("owner removes a member", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		teamRepo.EXPECT().GetMemberByUserID(context.Background(), uint(1), uint(7)).Return(&model.TeamMember{TeamID: 1, UserID: 7, Role: "owner"}, nil)
		teamRepo.EXPECT().DeleteMember(context.Background(), uint(1), uint(5)).Return(nil)

		if err := svc.RemoveTeamMember(context.Background(), 1, 5, 7, false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestTeamService_GetTeamByID(t *testing.T) {
	t.Run("returns ErrNotFound when team missing", func(t *testing.T) {
		svc, teamRepo, _, _, _ := newTestTeamService(t)
		teamRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(nil, nil)

		_, err := svc.GetTeamByID(context.Background(), 1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("aggregates member/project/app counts", func(t *testing.T) {
		svc, teamRepo, projectRepo, appRepo, _ := newTestTeamService(t)
		team := &model.Team{ID: 1, Name: "team"}
		teamRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(team, nil)
		teamRepo.EXPECT().CountMembers(context.Background(), uint(1)).Return(int64(3), nil)
		projectRepo.EXPECT().CountByTeam(context.Background(), uint(1)).Return(int64(2), nil)
		appRepo.EXPECT().CountByTeam(context.Background(), uint(1)).Return(int64(5), nil)

		got, err := svc.GetTeamByID(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.MemberCount != 3 || got.ProjectCount != 2 || got.AppCount != 5 {
			t.Fatalf("unexpected stats: %+v", got)
		}
	})
}

func TestTeamService_GetUserStats(t *testing.T) {
	t.Run("returns ErrNotFound when user missing", func(t *testing.T) {
		svc, _, _, _, userRepo := newTestTeamService(t)
		userRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(nil, nil)

		_, err := svc.GetUserStats(context.Background(), 1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("skips memberships whose team can no longer be found", func(t *testing.T) {
		svc, teamRepo, _, _, userRepo := newTestTeamService(t)
		user := &model.User{ID: 1, Username: "alice"}
		userRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(user, nil)
		teamRepo.EXPECT().GetByUserID(context.Background(), uint(1)).Return([]model.TeamMember{
			{TeamID: 404, UserID: 1, Role: "member"},
		}, nil)
		teamRepo.EXPECT().GetByID(context.Background(), uint(404)).Return(nil, nil)

		got, err := svc.GetUserStats(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.TeamCount != 1 || len(got.Teams) != 0 {
			t.Fatalf("expected the missing team to be skipped from Teams but still counted, got %+v", got)
		}
	})

	t.Run("aggregates stats across memberships", func(t *testing.T) {
		svc, teamRepo, projectRepo, appRepo, userRepo := newTestTeamService(t)
		user := &model.User{ID: 1, Username: "alice"}
		team := &model.Team{ID: 5, Name: "platform"}
		userRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(user, nil)
		teamRepo.EXPECT().GetByUserID(context.Background(), uint(1)).Return([]model.TeamMember{
			{TeamID: 5, UserID: 1, Role: "owner"},
		}, nil)
		teamRepo.EXPECT().GetByID(context.Background(), uint(5)).Return(team, nil)
		teamRepo.EXPECT().CountMembers(context.Background(), uint(5)).Return(int64(4), nil)
		projectRepo.EXPECT().CountByTeam(context.Background(), uint(5)).Return(int64(2), nil)
		appRepo.EXPECT().CountByTeam(context.Background(), uint(5)).Return(int64(6), nil)

		got, err := svc.GetUserStats(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.TeamCount != 1 || got.ProjectCount != 2 || got.AppCount != 6 {
			t.Fatalf("unexpected totals: %+v", got)
		}
		if len(got.Teams) != 1 || got.Teams[0].Role != "owner" || got.Teams[0].MemberCount != 4 {
			t.Fatalf("unexpected team info: %+v", got.Teams)
		}
	})
}
