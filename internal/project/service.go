package project

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Service struct {
	repository   Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) GetAllProjects(ctx context.Context) ([]Project, error) {
	data, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetProjectByName(ctx context.Context, name string) (*Project, error) {
	data, err := s.repository.GetByName(name)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("project not found")
	}

	return data, nil
}

func (s *Service) SaveProject(ctx context.Context, req CreateProjectRequestDto) (*Project, error) {
	existing, err := s.repository.GetByName(req.Name)
	if err != gorm.ErrRecordNotFound {
		fmt.Println("sapo")
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("the project already exists")
	}

	data, err := s.repository.Save(&CreateProjectRecordDto{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
