package service

import (
	"context"
	"errors"
	"fmt"

	"git.4thena.io/4thena/abys/internal/constant"
	"git.4thena.io/4thena/abys/internal/integration/ci"
	"git.4thena.io/4thena/abys/internal/integration/forge"
	"git.4thena.io/4thena/abys/internal/model"
	"git.4thena.io/4thena/abys/internal/repository"
	"gorm.io/gorm"
)

type AppService struct {
	appRepository      repository.AppRepository
	projectRepository  repository.ProjectRepository
	templateRepository repository.TemplateRepository
	forge              forge.Forge
	ci                 ci.Ci
}

func NewAppService(appRepository repository.AppRepository, projectRepository repository.ProjectRepository, templateRepository repository.TemplateRepository, forge forge.Forge, ci ci.Ci) *AppService {
	return &AppService{
		appRepository,
		projectRepository,
		templateRepository,
		forge,
		ci,
	}
}

func NewAppServiceReadOnly(appRepository repository.AppRepository, projectRepository repository.ProjectRepository, templateRepository repository.TemplateRepository) *AppService {
	return &AppService{
		appRepository:      appRepository,
		projectRepository:  projectRepository,
		templateRepository: templateRepository,
	}
}

func (s *AppService) CreateApp(ctx context.Context, app *model.App) (*model.App, error) {
	project, err := s.projectRepository.GetByID(ctx, app.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("there was an error reading getting the project: %w", err)
	}
	if project == nil {
		return nil, fmt.Errorf("project doesn't exists")
	}

	template, err := s.templateRepository.GetByID(ctx, app.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("there was an error reading getting the template: %w", err)
	}
	if template == nil {
		return nil, fmt.Errorf("template doesn't exists")
	}

	existing, err := s.appRepository.GetByName(ctx, app.Name)
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("there was an error reading the data: %w", err)
	}
	if existing != nil {
		return nil, errors.New("app already exists")
	}

	repo, err := s.forge.CreateRepo(ctx, "4thena", app.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create the app: %w", err)
	}

	ciRepo, err := s.ci.ActivateRepo(ctx, repo.ID)
	if err != nil {
		s.forge.DeleteRepo(ctx, "4thena", app.Name)
		return nil, fmt.Errorf("failed to activate ci repo: %w", err)
	}

	app.RepoID = repo.ID
	app.RepoURL = repo.URL
	app.CloneURL = repo.CloneURL
	app.CiID = ciRepo.ID
	app.CiURL = ciRepo.URL
	app.Status = constant.StatusNew

	err = s.appRepository.Save(ctx, app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (s *AppService) GetAllApps(ctx context.Context) ([]model.App, error) {
	data, err := s.appRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *AppService) GetAppByID(ctx context.Context, id uint) (*model.App, error) {
	app, err := s.appRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, fmt.Errorf("app not found")
	}
	return app, nil
}

func (s *AppService) GetAppBuilds(ctx context.Context, id uint) ([]model.Build, error) {
	app, err := s.appRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, fmt.Errorf("app not found")
	}
	return s.ci.GetBuilds(ctx, app.CiID)
}

func (s *AppService) GetAppsByProject(ctx context.Context, id uint) ([]model.App, error) {
	project, err := s.projectRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	apps, err := s.appRepository.GetByProject(ctx, id)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *AppService) GetAppsByTemplate(ctx context.Context, id uint) ([]model.App, error) {
	template, err := s.templateRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, fmt.Errorf("template not found")
	}

	apps, err := s.appRepository.GetByTemplate(ctx, id)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *AppService) DeleteApp(ctx context.Context, id uint) error {
	app, err := s.appRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return fmt.Errorf("app doesn't exist")
	}
	return s.appRepository.Delete(ctx, app)
}
