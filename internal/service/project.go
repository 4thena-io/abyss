package service

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
)

type ProjectRepository interface {
	Save(ctx context.Context, project *model.Project) error
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

func (s *ProjectService) DeleteProject(ctx context.Context, id uint) error {
	project, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if project == nil {
		return ErrNotFound
	}
	return s.repository.Delete(ctx, project)
}
