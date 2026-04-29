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

func setup() *Server {
	db := database.Connection

	forgeProvider, err := forge.NewForge()
	if err != nil {
		log.Fatalf("failed to create forge provider: %s", err)
	}
	ciProvider, err := ci.NewCi()
	if err != nil {
		log.Fatalf("failed to create ci provider: %s", err)
	}
	gitClient := git.New(config.Environment.ForgeToken)

	appRepo := repository.NewAppRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	templateRepo := repository.NewTemplateRepository(db)
	deploymentRepo := repository.NewDeploymentRepository(db)
	teamRepo := repository.NewTeamRepository(db)

	appService := service.NewAppService(appRepo, projectRepo, templateRepo, forgeProvider, ciProvider, gitClient)
	deploymentService := service.NewDeploymentService(deploymentRepo)
	projectService := service.NewProjectService(projectRepo)
	templateService := service.NewTemplateService(templateRepo)
	repoService := service.NewRepoService(forgeProvider)
	teamService := service.NewTeamService(teamRepo, projectRepo, appRepo)

	r := router.New(
		handler.NewAppHandler(appService, deploymentService),
		handler.NewProjectHandler(projectService, appService),
		handler.NewTemplateHandler(templateService),
		handler.NewRepoHandler(repoService),
		handler.NewTeamHandler(teamService),
	)

	return newServer(Config{Host: "0.0.0.0", Port: "8000"}, r)
}
