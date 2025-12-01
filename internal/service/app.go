package service

import (
	"context"
	"errors"
	"fmt"

	"git.d4ramirez.com/project-abyss/abys-api/internal/dto"
	"git.d4ramirez.com/project-abyss/abys-api/internal/ci"
	"git.d4ramirez.com/project-abyss/abys-api/internal/forge"
	"git.d4ramirez.com/project-abyss/abys-api/internal/model"
	"git.d4ramirez.com/project-abyss/abys-api/internal/repository"
	"gorm.io/gorm"
)

type AppService struct {
	repository        repository.AppRepository
	forge      				forge.Forge
	ci								ci.Ci
}

func NewAppService(repository repository.AppRepository, forge forge.Forge, ci ci.Ci) *AppService {
	return &AppService{
		repository,
		forge,
		ci,
	}
}

func (s *AppService) GetAllApps(ctx context.Context) ([]model.App, error) {
	data, err := s.repository.GetAllApps(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *AppService) GetAppByName(ctx context.Context, name string) (*model.App, error) {
	return s.repository.GetAppByName(ctx, name)
}

func (s *AppService) CreateApp(ctx context.Context, request dto.CreateAppRequestDto) (*model.App, error) {
	existing, err := s.repository.GetAppByName(ctx, request.Name)
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("there was an error reading the data: %w", err)
	}
	if existing != nil {
		return nil, errors.New("app already exists")
	}

	gitResponse, err := s.forge.CreateRepo(ctx, "4thena", request.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create the app: %w", err)
	}

	ciResponse, err := s.ci.ActivateRepo(ctx, dto.ActivateRepoRequestDTO{
		Owner: "4thena",
		Name: request.Name,
		CloneUrl: gitResponse.CloneUrl,
		ForgeRemoteId: gitResponse.RepoId,
	})
	if err != nil {
		s.forge.DeleteRepo(ctx, "4thena", request.Name)
		return nil, fmt.Errorf("failed to activate ci repo: %w", err)
	}

	data := dto.CreateAppRecordDto{
		Name:        request.Name,
		Description: request.Description,
		RepoId:      gitResponse.RepoId,
		RepoUrl:     gitResponse.HtmlUrl,
		CloneUrl:    gitResponse.CloneUrl,
		CiId:        ciResponse.RepoId,
		CiUrl:       ciResponse.RepoUrl,
		Project:     request.Project,
	}

	return s.repository.SaveApp(ctx, &data)
}
