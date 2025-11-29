package app

import (
	"context"
	"errors"
	"fmt"

	"git.d4ramirez.com/project-abyss/abys-api/internal/config"
	"git.d4ramirez.com/project-abyss/abys-api/internal/integration/git"
	"gorm.io/gorm"
)

type Service struct {
	repository        Repository
	gitProvider      	git.Provider
}

func NewService(repository Repository, gitProvider git.Provider) *Service {
	return &Service{
		repository,
		gitProvider,
	}
}

func (service *Service) GetAllApps(ctx context.Context) ([]App, error) {
	data, err := service.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetAppByName(ctx context.Context, name string) (*App, error) {
	return s.repository.GetByName(name)
}

func (s *Service) CreateApp(ctx context.Context, request CreateAppRequestDto) (*App, error) {
	existing, err := s.repository.GetByName(request.Name)
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("there was an error reading the data: %w", err)
	}
	if existing != nil {
		return nil, errors.New("app already exists")
	}

	gitResponse, err := s.gitProvider.CreateRepo(ctx, "4thena", request.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create the app: %w", err)
	}

	data := CreateAppRecordDto{
		Name:        request.Name,
		Description: request.Description,
		RepoId:      gitResponse.RepoId,
		RepoUrl:     gitResponse.HtmlUrl,
		CloneUrl:    gitResponse.CloneUrl,
		CiId:        3,
		CiUrl:       config.Environment.WoodpeckerHost + "/repos/3",
		Project:     request.Project,
	}

	return s.repository.Save(&data)
}
