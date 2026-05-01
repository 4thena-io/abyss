package main

import (
	"log"

	"github.com/4thena-io/abyss/internal/api/rest/handler"
	"github.com/4thena-io/abyss/internal/api/rest/router"
	"github.com/4thena-io/abyss/internal/config"
	"github.com/4thena-io/abyss/internal/database"
	"github.com/4thena-io/abyss/internal/integration/ci"
	"github.com/4thena-io/abyss/internal/integration/forge"
	"github.com/4thena-io/abyss/internal/integration/git"
	"github.com/4thena-io/abyss/internal/repository"
	"github.com/4thena-io/abyss/internal/service"
)

func setup(cfg *config.Config) *Server {
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	forgeProvider, err := forge.NewForge(cfg.Forge)
	if err != nil {
		log.Fatalf("failed to create forge provider: %v", err)
	}
	ciProvider, err := ci.NewCi(cfg.CI, cfg.Forge)
	if err != nil {
		log.Fatalf("failed to create ci provider: %v", err)
	}
	gitClient := git.New(cfg.Forge.Token)

	appRepo := repository.NewAppRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	deploymentRepo := repository.NewDeploymentRepository(db)
	teamRepo := repository.NewTeamRepository(db)

	appService := service.NewAppService(appRepo, projectRepo, templateRepo, forgeProvider, ciProvider, gitClient, cfg.Forge.Owner)
	deploymentService := service.NewDeploymentService(deploymentRepo)
	projectService := service.NewProjectService(projectRepo)
	templateService := service.NewTemplateService(templateRepo)
	repoService := service.NewRepoService(forgeProvider, cfg.Forge.Owner)
	teamService := service.NewTeamService(teamRepo, projectRepo, appRepo)

	r := router.New(
		handler.NewAppHandler(appService, deploymentService),
		handler.NewProjectHandler(projectService, appService),
		handler.NewTemplateHandler(templateService),
		handler.NewRepoHandler(repoService),
		handler.NewTeamHandler(teamService),
	)

	return newServer(Config{Host: cfg.Server.Host, Port: cfg.Server.Port}, r)
}
