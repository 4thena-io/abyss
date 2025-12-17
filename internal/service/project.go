package service

import (
	"context"
	"fmt"

	"git.4thena.io/4thena/abys/internal/dto"
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

func (s *ProjectService) GetProjectByName(ctx context.Context, name string) (*model.Project, error) {
	data, err := s.repository.GetProjectByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("project not found")
	}

	return data, nil
}

func (s *ProjectService) SaveProject(ctx context.Context, req dto.CreateProjectRequestDto) (*model.Project, error) {
	existing, err := s.repository.GetProjectByName(ctx, req.Name)
	if err != gorm.ErrRecordNotFound {
		fmt.Println("sapo")
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("the project already exists")
	}

	data, err := s.repository.SaveProject(ctx, &dto.CreateProjectRecordDto{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
