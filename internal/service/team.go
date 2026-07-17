package service

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
)

type TeamRepository interface {
	Save(ctx context.Context, team *model.Team) error
	Update(ctx context.Context, team *model.Team) error
	GetAll(ctx context.Context) ([]model.Team, error)
	GetByID(ctx context.Context, id uint) (*model.Team, error)
	GetByName(ctx context.Context, name string) (*model.Team, error)
	GetMembers(ctx context.Context, teamID uint) ([]model.TeamMember, error)
	GetByUserID(ctx context.Context, userID uint) ([]model.TeamMember, error)
	SaveMember(ctx context.Context, member *model.TeamMember) error
	DeleteMember(ctx context.Context, teamID, memberID uint) error
	CountMembers(ctx context.Context, teamID uint) (int64, error)
	Delete(ctx context.Context, team *model.Team) error
}

type UserLookup interface {
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

type TeamStats struct {
	Team         model.Team
	MemberCount  int64
	ProjectCount int64
	AppCount     int64
}

type UserTeamInfo struct {
	Team        model.Team
	Role        string
	MemberCount int64
	ProjectCount int64
}

type UserStats struct {
	User         model.User
	TeamCount    int64
	ProjectCount int64
	AppCount     int64
	Teams        []UserTeamInfo
}

type TeamService struct {
	teamRepo    TeamRepository
	projectRepo ProjectRepository
	appRepo     AppRepository
	userRepo    UserLookup
}

func NewTeamService(teamRepo TeamRepository, projectRepo ProjectRepository, appRepo AppRepository, userRepo UserLookup) *TeamService {
	return &TeamService{teamRepo, projectRepo, appRepo, userRepo}
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

func (s *TeamService) SaveTeam(ctx context.Context, team *model.Team, creatorID uint) (*model.Team, error) {
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

	if err := s.teamRepo.SaveMember(ctx, &model.TeamMember{
		TeamID: team.ID,
		UserID: creatorID,
		Role:   "owner",
	}); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) UpdateTeam(ctx context.Context, id uint, name, description string) (*model.Team, error) {
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, ErrNotFound
	}

	team.Name = name
	team.Description = description

	if err := s.teamRepo.Update(ctx, team); err != nil {
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

func (s *TeamService) AddTeamMember(ctx context.Context, teamID uint, username, role string) (*model.TeamMember, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	member := &model.TeamMember{
		TeamID: teamID,
		UserID: user.ID,
		Role:   role,
	}
	if err := s.teamRepo.SaveMember(ctx, member); err != nil {
		return nil, err
	}
	member.User = *user
	return member, nil
}

func (s *TeamService) RemoveTeamMember(ctx context.Context, teamID, memberID uint) error {
	return s.teamRepo.DeleteMember(ctx, teamID, memberID)
}

func (s *TeamService) GetUserStats(ctx context.Context, userID uint) (*UserStats, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	memberships, err := s.teamRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	teams := make([]UserTeamInfo, 0, len(memberships))
	var totalProjects, totalApps int64
	for _, m := range memberships {
		team, err := s.teamRepo.GetByID(ctx, m.TeamID)
		if err != nil || team == nil {
			continue
		}
		memberCount, _ := s.teamRepo.CountMembers(ctx, team.ID)
		projectCount, _ := s.projectRepo.CountByTeam(ctx, team.ID)
		appCount, _ := s.appRepo.CountByTeam(ctx, team.ID)
		totalProjects += projectCount
		totalApps += appCount
		teams = append(teams, UserTeamInfo{
			Team:         *team,
			Role:         m.Role,
			MemberCount:  memberCount,
			ProjectCount: projectCount,
		})
	}

	return &UserStats{
		User:         *user,
		TeamCount:    int64(len(memberships)),
		ProjectCount: totalProjects,
		AppCount:     totalApps,
		Teams:        teams,
	}, nil
}
