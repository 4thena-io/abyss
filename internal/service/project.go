package service

import (
	"context"
	"fmt"

	"git.4thena.io/4thena/abys/internal/model"
	"git.4thena.io/4thena/abys/internal/repository"
	"gorm.io/gorm"
)

type ProjectService struct {
	repository   repository.ProjectRepository
}

func NewProjectService(repository repository.ProjectRepository) *ProjectService {
	return &ProjectService{
		repository,
	}
}

func (s *ProjectService) GetAllProjects(ctx context.Context) ([]model.Project, error) {
	data, err := s.repository.GetAllProjects(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id uint) (*model.Project, error) {
	data, err := s.repository.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("project not found")
	}

	return data, nil
}

func (s *ProjectService) SaveProject(ctx context.Context, project *model.Project) (*model.Project, error) {
	existing, err := s.repository.GetProjectByName(ctx, project.Name)
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("the project already exists")
	}

 	err = s.repository.SaveProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}
