package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/4thena-io/abyss/internal/api/rest/handler"
	"github.com/4thena-io/abyss/internal/api/rest/middleware"
	"github.com/4thena-io/abyss/internal/api/rest/router"
	"github.com/4thena-io/abyss/internal/config"
	"github.com/4thena-io/abyss/internal/database"
	"github.com/4thena-io/abyss/internal/git"
	"github.com/4thena-io/abyss/internal/integration/ci"
	"github.com/4thena-io/abyss/internal/integration/forge"
	"github.com/4thena-io/abyss/internal/repository"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func setup(cfg *config.Config, configPath string, restartCh chan<- struct{}) *Server {
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	if err := database.RunMigrations(db); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	setupHandler := handler.NewSetupHandler(configPath, cfg.Auth.ClientID != "", restartCh)

	// When OAuth is not configured yet, skip forge/CI setup — the setup wizard handles it.
	var (
		gitClient     *git.GitClient
		authService   *service.AuthService
		baseURL       string
		webhookSecret string
	)

	appRepo := repository.NewAppRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	deploymentRepo := repository.NewDeploymentRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)

	branch := cfg.Forge.Branch
	if branch == "" {
		branch = "main"
	}

	appService := service.NewAppService(appRepo, projectRepo, templateRepo, nil, nil, gitClient, cfg.Forge.Owner, "", "", branch)
	deploymentService := service.NewDeploymentService(deploymentRepo)
	projectService := service.NewProjectService(projectRepo)
	templateService := service.NewTemplateService(templateRepo)
	repoService := service.NewRepoService(nil, cfg.Forge.Owner)
	teamService := service.NewTeamService(teamRepo, projectRepo, appRepo, userRepo)
	docsService := service.NewDocsService(appRepo, gitClient, filepath.Join(config.DefaultDataDir(), "docs"))

	if cfg.Auth.ClientID != "" {
		forgeProvider, err := forge.NewForge(cfg.Forge)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create forge provider")
		}
		ciProvider, err := ci.NewCi(cfg.CI, cfg.Forge)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create ci provider")
		}

		botUser, err := forgeProvider.GetAuthenticatedUser(context.Background())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get forge bot user")
		}
		log.Info().Str("username", botUser.Username).Msg("forge bot user identified")

		jwtSecret, err := ensureJWTSecret(db)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to provision jwt secret")
		}

		webhookSecret, err = ensureWebhookSecret(db)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to provision webhook secret")
		}

		host := cfg.Server.Host
		if host == "" || host == "0.0.0.0" {
			host = "localhost"
		}
		baseURL = fmt.Sprintf("http://%s:%s", host, cfg.Server.Port)

		gitClient = git.New(cfg.Forge.Token, botUser.FullName, botUser.Email)

		appService = service.NewAppService(appRepo, projectRepo, templateRepo, forgeProvider, ciProvider, gitClient, cfg.Forge.Owner, baseURL, webhookSecret, branch)
		repoService = service.NewRepoService(forgeProvider, cfg.Forge.Owner)
		docsService = service.NewDocsService(appRepo, gitClient, filepath.Join(config.DefaultDataDir(), "docs"))

		authService = service.NewAuthService(
			userRepo,
			forgeProvider,
			cfg.Auth.ClientID,
			cfg.Auth.ClientSecret,
			cfg.Forge.Type,
			cfg.Forge.Host,
			cfg.Auth.CallbackURL,
			jwtSecret,
			cfg.Forge.Owner,
		)
	}

	requireSetup := middleware.RequireSetup(cfg.Auth.ClientID)

	// Pass a nil interface (not a nil *service.AuthService) so RequireAuth can detect unconfigured state.
	var authProvider middleware.AuthProvider
	if authService != nil {
		authProvider = authService
	}
	requireAuth := middleware.RequireAuth(authProvider)

	r := router.New(
		requireSetup,
		requireAuth,
		setupHandler,
		handler.NewAuthHandler(authService),
		handler.NewAppHandler(appService, deploymentService),
		handler.NewProjectHandler(projectService, appService),
		handler.NewTemplateHandler(templateService),
		handler.NewRepoHandler(repoService),
		handler.NewTeamHandler(teamService),
		handler.NewUserHandler(teamService),
		handler.NewHookHandler(docsService, branch, webhookSecret),
		handler.NewDocsHandler(docsService),
	)

	return newServer(Config{Host: cfg.Server.Host, Port: cfg.Server.Port}, r)
}

func ensureWebhookSecret(db *gorm.DB) (string, error) {
	return ensureSecret(db, "webhook-secret")
}

// ensureJWTSecret loads the JWT signing secret from the database, generating and
// persisting a new one if it does not exist yet. This keeps the secret out of
// config files and rotates automatically on a fresh install.
func ensureJWTSecret(db *gorm.DB) (string, error) {
	return ensureSecret(db, "jwt-secret")
}

func ensureSecret(db *gorm.DB, key string) (string, error) {
	repo := repository.NewServerConfigRepository(db)
	ctx := context.Background()

	secret, err := repo.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", key, err)
	}
	if secret != "" {
		return secret, nil
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate %s: %w", key, err)
	}
	secret = hex.EncodeToString(b)

	if err := repo.Set(ctx, key, secret); err != nil {
		return "", fmt.Errorf("store %s: %w", key, err)
	}
	log.Info().Str("key", key).Msg("generated new secret and stored it in database")
	return secret, nil
}
