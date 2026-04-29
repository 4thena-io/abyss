package service

import (
	"context"
	"fmt"

	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/repository"
)

type TeamStats struct {
	Team         model.Team
	MemberCount  int64
	ProjectCount int64
	AppCount     int64
}

type TeamService struct {
	teamRepo    repository.TeamRepository
	projectRepo repository.ProjectRepository
	appRepo     repository.AppRepository
}

func NewTeamService(teamRepo repository.TeamRepository, projectRepo repository.ProjectRepository, appRepo repository.AppRepository) *TeamService {
	return &TeamService{teamRepo, projectRepo, appRepo}
}

func (s *TeamService) GetAllTeams(ctx context.Context) ([]TeamStats, error) {
	teams, err := s.teamRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]TeamStats, len(teams))
	for i, team := range teams {
		memberCount, _ := s.teamRepo.CountMembers(ctx, team.ID)
		projectCount, _ := s.projectRepo.CountByTeam(ctx, team.ID)
		appCount, _ := s.appRepo.CountByTeam(ctx, team.ID)
		result[i] = TeamStats{
			Team:         team,
			MemberCount:  memberCount,
			ProjectCount: projectCount,
			AppCount:     appCount,
		}
	}
	return result, nil
}

func (s *TeamService) GetTeamByID(ctx context.Context, id uint) (*TeamStats, error) {
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, fmt.Errorf("team not found")
	}

	memberCount, _ := s.teamRepo.CountMembers(ctx, id)
	projectCount, _ := s.projectRepo.CountByTeam(ctx, id)
	appCount, _ := s.appRepo.CountByTeam(ctx, id)

	return &TeamStats{
		Team:         *team,
		MemberCount:  memberCount,
		ProjectCount: projectCount,
		AppCount:     appCount,
	}, nil
}

func (s *TeamService) SaveTeam(ctx context.Context, team *model.Team) (*model.Team, error) {
	existing, err := s.teamRepo.GetByName(ctx, team.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("team already exists")
	}

	if err := s.teamRepo.Save(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, id uint) error {
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if team == nil {
		return fmt.Errorf("team not found")
	}
	return s.teamRepo.Delete(ctx, team)
}

func (s *TeamService) GetTeamProjects(ctx context.Context, teamID uint) ([]model.Project, error) {
	return s.projectRepo.GetByTeam(ctx, teamID)
}

func (s *TeamService) GetTeamMembers(ctx context.Context, teamID uint) ([]model.TeamMember, error) {
	return s.teamRepo.GetMembers(ctx, teamID)
}

func (s *TeamService) AddTeamMember(ctx context.Context, member *model.TeamMember) (*model.TeamMember, error) {
	if err := s.teamRepo.SaveMember(ctx, member); err != nil {
		return nil, err
	}
	return member, nil
}
