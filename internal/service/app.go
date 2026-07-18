package service

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/4thena-io/abyss/internal/git"
	"github.com/4thena-io/abyss/internal/integration/ci"
	"github.com/4thena-io/abyss/internal/integration/forge"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/rs/zerolog/log"
)

type AppRepository interface {
	Save(ctx context.Context, app *model.App) error
	GetAll(ctx context.Context) ([]model.App, error)
	GetByID(ctx context.Context, id uint) (*model.App, error)
	GetByName(ctx context.Context, name string) (*model.App, error)
	GetByProject(ctx context.Context, id uint) ([]model.App, error)
	GetByTemplate(ctx context.Context, id uint) ([]model.App, error)
	Update(ctx context.Context, app *model.App) error
	Delete(ctx context.Context, app *model.App) error
	CountByTeam(ctx context.Context, teamID uint) (int64, error)
}

type AppService struct {
	appRepository      AppRepository
	projectRepository  ProjectRepository
	templateRepository TemplateRepository
	forge              forge.Forge
	ci                 ci.CI
	git                *git.GitClient
	owner              string
	baseURL            string
	webhookSecret      string
	branch             string
}

func NewAppService(
	appRepository AppRepository,
	projectRepository ProjectRepository,
	templateRepository TemplateRepository,
	forge forge.Forge,
	ci ci.CI,
	gitClient *git.GitClient,
	owner string,
	baseURL string,
	webhookSecret string,
	branch string,
) *AppService {
	return &AppService{
		appRepository:      appRepository,
		projectRepository:  projectRepository,
		templateRepository: templateRepository,
		forge:              forge,
		ci:                 ci,
		git:                gitClient,
		owner:              owner,
		baseURL:            baseURL,
		webhookSecret:      webhookSecret,
		branch:             branch,
	}
}

// registerDocsWebhook creates a push-event webhook on the repo so that changes
// to docs/ or .abyss.yml trigger a re-render. Failures are non-fatal.
func (s *AppService) registerDocsWebhook(ctx context.Context, app *model.App) {
	if s.baseURL == "" {
		return
	}
	owner, repoName, ok := strings.Cut(app.RepoFullName, "/")
	if !ok {
		return
	}
	callbackURL := fmt.Sprintf("%s/api/hooks/forge/%d?access_token=%s", s.baseURL, app.ID, url.QueryEscape(s.webhookSecret))
	if err := s.forge.CreateWebhook(ctx, owner, repoName, callbackURL, s.webhookSecret, s.branch); err != nil {
		log.Warn().Err(err).Uint("app_id", app.ID).Msg("failed to register docs webhook")
	}
}

func (s *AppService) CreateAppFromTemplate(ctx context.Context, app *model.App) (*model.App, error) {
	// ProjectID is optional — apps can be standalone (e.g. shared libraries).
	if app.ProjectID != nil {
		project, err := s.projectRepository.GetByID(ctx, *app.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("there was an error reading getting the project: %w", err)
		}
		if project == nil {
			return nil, fmt.Errorf("project doesn't exists")
		}
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
	if err != nil {
		return nil, fmt.Errorf("there was an error reading the data: %w", err)
	}
	if existing != nil {
		return nil, ErrConflict
	}

	// Create empty repo in forge
	repo, err := s.forge.CreateRepo(ctx, s.owner, app.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create the app: %w", err)
	}
	// Clone template, prepare files, and push to new repo
	if err := s.initRepoFromTemplate(template.CloneURL, repo.CloneURL, app.Name); err != nil {
		if delErr := s.forge.DeleteRepo(ctx, s.owner, app.Name); delErr != nil {
			log.Warn().Err(delErr).Str("repo", app.Name).Msg("failed to roll back forge repo after template init failure")
		}
		return nil, fmt.Errorf("failed to initialize repo from template: %w", err)
	}

	// Activate CI for this repo
	ciRepo, err := s.ci.ActivateRepo(ctx, repo.ID, repo.FullName)
	if err != nil {
		if delErr := s.forge.DeleteRepo(ctx, s.owner, app.Name); delErr != nil {
			log.Warn().Err(delErr).Str("repo", app.Name).Msg("failed to roll back forge repo after ci activation failure")
		}
		return nil, fmt.Errorf("failed to activate ci repo: %w", err)
	}

	app.RepoID = repo.ID
	app.RepoFullName = repo.FullName
	app.RepoURL = repo.URL
	app.CloneURL = repo.CloneURL
	app.CIID = ciRepo.ID
	app.CISlug = ciRepo.Slug
	app.CIURL = ciRepo.URL


	err = s.appRepository.Save(ctx, app)
	if err != nil {
		return nil, err
	}

	s.registerDocsWebhook(ctx, app)
	return app, nil
}

// initRepoFromTemplate clones template, removes .git, ensures .abyss.yml, and pushes to new repo.
func (s *AppService) initRepoFromTemplate(templateCloneURL, newRepoCloneURL, appName string) error {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "abyss-app-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

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
	if app.ProjectID != nil {
		project, err := s.projectRepository.GetByID(ctx, *app.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("there was an error reading getting the project: %w", err)
		}
		if project == nil {
			return nil, fmt.Errorf("project doesn't exists")
		}
	}

	// Validate app name uniqueness
	existing, err := s.appRepository.GetByName(ctx, app.Name)
	if err != nil {
		return nil, fmt.Errorf("there was an error reading the data: %w", err)
	}
	if existing != nil {
		return nil, ErrConflict
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


	err = s.appRepository.Save(ctx, app)
	if err != nil {
		return nil, err
	}

	s.registerDocsWebhook(ctx, app)
	return app, nil
}

// ensureAbyssConfig clones the repo, checks for .abyss.yml, creates if missing, and pushes.
func (s *AppService) ensureAbyssConfig(repoCloneURL, appName string) error {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "abyss-app-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

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
	return s.appRepository.GetByID(ctx, id)
}

func (s *AppService) GetAppBuilds(ctx context.Context, id uint) ([]model.Build, error) {
	app, err := s.appRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, ErrNotFound
	}
	return s.ci.GetBuilds(ctx, app.CIID, app.CISlug)
}

func (s *AppService) GetAppsByProject(ctx context.Context, id uint) ([]model.App, error) {
	project, err := s.projectRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrNotFound
	}

	return s.appRepository.GetByProject(ctx, id)
}

func (s *AppService) GetAppsByTemplate(ctx context.Context, id uint) ([]model.App, error) {
	template, err := s.templateRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, ErrNotFound
	}

	return s.appRepository.GetByTemplate(ctx, id)
}

func (s *AppService) RepairWebhook(ctx context.Context, id uint) error {
	app, err := s.appRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return ErrNotFound
	}

	owner, repoName, ok := strings.Cut(app.RepoFullName, "/")
	if !ok {
		return fmt.Errorf("invalid repo full name: %s", app.RepoFullName)
	}

	callbackURL := fmt.Sprintf("%s/api/hooks/forge/%d?access_token=%s", s.baseURL, app.ID, url.QueryEscape(s.webhookSecret))

	if err := s.forge.DeleteWebhook(ctx, owner, repoName, callbackURL); err != nil {
		return fmt.Errorf("failed to remove old webhook: %w", err)
	}
	if err := s.forge.CreateWebhook(ctx, owner, repoName, callbackURL, s.webhookSecret, s.branch); err != nil {
		return fmt.Errorf("failed to register webhook: %w", err)
	}
	return nil
}

func (s *AppService) UpdateApp(ctx context.Context, id, callerID uint, isAdmin bool, name, description string, projectID *uint) (*model.App, error) {
	app, err := s.appRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, ErrNotFound
	}
	if !isAdmin && app.CreatorID != callerID {
		return nil, ErrForbidden
	}

	app.Name = name
	app.Description = description
	app.ProjectID = projectID

	if err := s.appRepository.Update(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *AppService) DeleteApp(ctx context.Context, id, callerID uint, isAdmin bool) error {
	app, err := s.appRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return ErrNotFound
	}
	if !isAdmin && app.CreatorID != callerID {
		return ErrForbidden
	}

	err = s.ci.DeleteRepo(ctx, app.CIID, app.CISlug)
	if err != nil {
		return err
	}

	err = s.forge.DeleteRepo(ctx, s.owner, app.Name)
	if err != nil {
		return err
	}

	return s.appRepository.Delete(ctx, app)
}
