package service

import (
	"context"
	"errors"
	"fmt"

	"git.4thena.io/4thena/abys/internal/ci"
	"git.4thena.io/4thena/abys/internal/constant"
	"git.4thena.io/4thena/abys/internal/dto/request"
	"git.4thena.io/4thena/abys/internal/dto/response"
	"git.4thena.io/4thena/abys/internal/forge"
	"git.4thena.io/4thena/abys/internal/model"
	"git.4thena.io/4thena/abys/internal/repository"
	"gorm.io/gorm"
)

type AppService struct {
	repository repository.AppRepository
	forge      forge.Forge
	ci         ci.Ci
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

func (s *AppService) GetAppByID(ctx context.Context, id uint) (*model.App, error) {
	return s.repository.GetAppById(ctx, id)
}

func (s *AppService) CreateApp(ctx context.Context, app *model.App) (*model.App, error) {
	existing, err := s.repository.GetAppByName(ctx, app.Name)
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("there was an error reading the data: %w", err)
	}
	if existing != nil {
		return nil, errors.New("app already exists")
	}

	gitResponse, err := s.forge.CreateRepo(ctx, "4thena", app.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create the app: %w", err)
	}

	ciResponse, err := s.ci.ActivateRepo(ctx, request.ActivateRepo{
		Owner:         "4thena",
		Name:          app.Name,
		CloneUrl:      gitResponse.CloneUrl,
		ForgeRemoteId: gitResponse.RepoId,
	})
	if err != nil {
		s.forge.DeleteRepo(ctx, "4thena", app.Name)
		return nil, fmt.Errorf("failed to activate ci repo: %w", err)
	}

	app.RepoID = gitResponse.RepoId
	app.RepoURL = gitResponse.HtmlUrl
	app.CloneURL = gitResponse.CloneUrl
	app.CiID = ciResponse.RepoId
	app.CiURL = ciResponse.RepoUrl
	app.Status = constant.StatusNew

	err = s.repository.SaveApp(ctx, app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (s *AppService) GetAppBuilds(ctx context.Context, appID uint) ([]response.Build, error) {
	app, err := s.repository.GetAppById(ctx, appID)
	if err != nil {
		return nil, err
	}
	return s.ci.GetBuilds(ctx, app.CiID)
}
