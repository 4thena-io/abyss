package service

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
)

type TeamRepository interface {
	Save(ctx context.Context, team *model.Team) error
	GetAll(ctx context.Context) ([]model.Team, error)
	GetByID(ctx context.Context, id uint) (*model.Team, error)
	GetByName(ctx context.Context, name string) (*model.Team, error)
	GetMembers(ctx context.Context, teamID uint) ([]model.TeamMember, error)
	SaveMember(ctx context.Context, member *model.TeamMember) error
	CountMembers(ctx context.Context, teamID uint) (int64, error)
	Delete(ctx context.Context, team *model.Team) error
}

type TeamStats struct {
	Team         model.Team
	MemberCount  int64
	ProjectCount int64
	AppCount     int64
}

type TeamService struct {
	teamRepo    TeamRepository
	projectRepo ProjectRepository
	appRepo     AppRepository
}

func NewTeamService(teamRepo TeamRepository, projectRepo ProjectRepository, appRepo AppRepository) *TeamService {
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
		return nil, ErrNotFound
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
		return nil, ErrConflict
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
		return ErrNotFound
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
