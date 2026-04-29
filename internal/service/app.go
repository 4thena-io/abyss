package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/4thena-io/abyss/internal/constant"
	"github.com/4thena-io/abyss/internal/integration/ci"
	"github.com/4thena-io/abyss/internal/integration/forge"
	"github.com/4thena-io/abyss/internal/integration/git"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/repository"
	"gorm.io/gorm"
)

type AppService struct {
	appRepository      repository.AppRepository
	projectRepository  repository.ProjectRepository
	templateRepository repository.TemplateRepository
	forge              forge.Forge
	ci                 ci.CI
	git                *git.GitClient
}

func NewAppService(
	appRepository repository.AppRepository,
	projectRepository repository.ProjectRepository,
	templateRepository repository.TemplateRepository,
	forge forge.Forge,
	ci ci.CI,
	gitClient *git.GitClient,
) *AppService {
	return &AppService{
		appRepository:      appRepository,
		projectRepository:  projectRepository,
		templateRepository: templateRepository,
		forge:              forge,
		ci:                 ci,
		git:                gitClient,
	}
}

func (s *AppService) CreateAppFromTemplate(ctx context.Context, app *model.App) (*model.App, error) {
	project, err := s.projectRepository.GetByID(ctx, app.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("there was an error reading getting the project: %w", err)
	}
	if project == nil {
		return nil, fmt.Errorf("project doesn't exists")
	}

	if app.TemplateID == nil {
		return nil, fmt.Errorf("template ID is required")
	}

	template, err := s.templateRepository.GetByID(ctx, *app.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("there was an error reading getting the template: %w", err)
	}
	if template == nil {
		return nil, fmt.Errorf("template doesn't exists")
	}

	existing, err := s.appRepository.GetByName(ctx, app.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("there was an error reading the data: %w", err)
	}
	if existing != nil {
		return nil, errors.New("app already exists")
	}

	// Create empty repo in forge
	repo, err := s.forge.CreateRepo(ctx, "4thena", app.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create the app: %w", err)
	}
	// Clone template, prepare files, and push to new repo
	if err := s.initRepoFromTemplate(template.CloneURL, repo.CloneURL, app.Name); err != nil {
		s.forge.DeleteRepo(ctx, "4thena", app.Name)
		return nil, fmt.Errorf("failed to initialize repo from template: %w", err)
	}

	// Activate CI for this repo
	ciRepo, err := s.ci.ActivateRepo(ctx, repo.ID, repo.FullName)
	if err != nil {
		s.forge.DeleteRepo(ctx, "4thena", app.Name)
		return nil, fmt.Errorf("failed to activate ci repo: %w", err)
	}

	app.RepoID = repo.ID
	app.RepoFullName = repo.FullName
	app.RepoURL = repo.URL
	app.CloneURL = repo.CloneURL
	app.CIID = ciRepo.ID
	app.CISlug = ciRepo.Slug
	app.CIURL = ciRepo.URL
	app.Status = constant.StatusNew

	err = s.appRepository.Save(ctx, app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// initRepoFromTemplate clones template, removes .git, ensures .abyss.yml, and pushes to new repo.
func (s *AppService) initRepoFromTemplate(templateCloneURL, newRepoCloneURL, appName string) error {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "abyss-app-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	localPath := filepath.Join(tmpDir, "repo")

	// Clone template repo
	if err := s.git.Clone(templateCloneURL, localPath); err != nil {
		return fmt.Errorf("failed to clone template: %w", err)
	}

	// Remove .git directory
	if err := s.git.RemoveGitDir(localPath); err != nil {
		return fmt.Errorf("failed to remove .git: %w", err)
	}

	// Ensure .abyss.yml exists
	if !s.git.FileExists(localPath, ".abyss.yml") {
		content := s.generateAbyssConfig(appName)
		if err := s.git.CreateFile(localPath, ".abyss.yml", content); err != nil {
			return fmt.Errorf("failed to create .abyss.yml: %w", err)
		}
	}

	// Init new repo and push to remote
	if err := s.git.InitAndPush(localPath, newRepoCloneURL, "Initial commit from template"); err != nil {
		return fmt.Errorf("failed to push to new repo: %w", err)
	}

	return nil
}

func (s *AppService) CreateAppFromRepo(ctx context.Context, app *model.App) (*model.App, error) {
	// Validate project exists (if provided)
	if app.ProjectID != 0 {
		project, err := s.projectRepository.GetByID(ctx, app.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("there was an error reading getting the project: %w", err)
		}
		if project == nil {
			return nil, fmt.Errorf("project doesn't exists")
		}
	}

	// Validate app name uniqueness
	existing, err := s.appRepository.GetByName(ctx, app.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("there was an error reading the data: %w", err)
	}
	if existing != nil {
		return nil, errors.New("app already exists")
	}

	// Validate repo exists in forge
	repo, err := s.forge.GetRepo(ctx, app.RepoID)
	if err != nil {
		return nil, fmt.Errorf("repo doesn't exist: %w", err)
	}

	// Check for .abyss.yml, create if missing
	if err := s.ensureAbyssConfig(repo.CloneURL, app.Name); err != nil {
		return nil, fmt.Errorf("failed to ensure .abyss.yml: %w", err)
	}

	// Activate CI for this repo
	ciRepo, err := s.ci.ActivateRepo(ctx, repo.ID, repo.FullName)
	if err != nil {
		return nil, fmt.Errorf("failed to activate ci repo: %w", err)
	}

	app.RepoID = repo.ID
	app.RepoFullName = repo.FullName
	app.RepoURL = repo.URL
	app.CloneURL = repo.CloneURL
	app.CIID = ciRepo.ID
	app.CISlug = ciRepo.Slug
	app.CIURL = ciRepo.URL
	app.Status = constant.StatusNew

	err = s.appRepository.Save(ctx, app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// ensureAbyssConfig clones the repo, checks for .abyss.yml, creates if missing, and pushes.
func (s *AppService) ensureAbyssConfig(repoCloneURL, appName string) error {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "abyss-app-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	localPath := filepath.Join(tmpDir, "repo")

	// Clone the repo
	if err := s.git.Clone(repoCloneURL, localPath); err != nil {
		return fmt.Errorf("failed to clone repo: %w", err)
	}

	// Check if .abyss.yml exists
	if s.git.FileExists(localPath, ".abyss.yml") {
		return nil // Already exists, nothing to do
	}

	// Create .abyss.yml
	content := s.generateAbyssConfig(appName)
	if err := s.git.CreateFile(localPath, ".abyss.yml", content); err != nil {
		return fmt.Errorf("failed to create .abyss.yml: %w", err)
	}

	// Commit and push
	if err := s.git.AddCommitPush(localPath, "Initialize .abyss.yml"); err != nil {
		return fmt.Errorf("failed to push .abyss.yml: %w", err)
	}

	return nil
}

func (s *AppService) generateAbyssConfig(appName string) string {
	return fmt.Sprintf(`name: %s
description: ""

dependencies:
  apps: []
  services: []

docs:
  folder: ./docs
`, appName)
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
	return s.ci.GetBuilds(ctx, app.CIID, app.CISlug)
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

	err = s.ci.DeleteRepo(ctx, app.CIID, app.CISlug)
	if err != nil {
		return err
	}

	err = s.forge.DeleteRepo(ctx, "4thena", app.Name)
	if err != nil {
		return err
	}

	return s.appRepository.Delete(ctx, app)
}
