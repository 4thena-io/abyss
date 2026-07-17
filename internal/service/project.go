package service

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
)

type ProjectRepository interface {
	Save(ctx context.Context, project *model.Project) error
	Update(ctx context.Context, project *model.Project) error
	GetAll(ctx context.Context) ([]model.Project, error)
	GetByID(ctx context.Context, id uint) (*model.Project, error)
	GetByName(ctx context.Context, name string) (*model.Project, error)
	GetByTeam(ctx context.Context, teamID uint) ([]model.Project, error)
	Delete(ctx context.Context, project *model.Project) error
	CountByTeam(ctx context.Context, teamID uint) (int64, error)
}

type ProjectService struct {
	repository ProjectRepository
}

func NewProjectService(repository ProjectRepository) *ProjectService {
	return &ProjectService{
		repository,
	}
}

func (s *ProjectService) GetAllProjects(ctx context.Context) ([]model.Project, error) {
	data, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id uint) (*model.Project, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *ProjectService) SaveProject(ctx context.Context, project *model.Project) (*model.Project, error) {
	existing, err := s.repository.GetByName(ctx, project.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrConflict
	}

	err = s.repository.Save(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id, callerID uint, isAdmin bool, name, description string, teamID *uint) (*model.Project, error) {
	project, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrNotFound
	}
	if !isAdmin && project.CreatorID != callerID {
		return nil, ErrForbidden
	}

	project.Name = name
	project.Description = description
	project.TeamID = teamID

	if err := s.repository.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, id, callerID uint, isAdmin bool) error {
	project, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if project == nil {
		return ErrNotFound
	}
	if !isAdmin && project.CreatorID != callerID {
		return ErrForbidden
	}
	return s.repository.Delete(ctx, project)
}
