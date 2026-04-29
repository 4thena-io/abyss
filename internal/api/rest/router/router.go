package router

import (
	"github.com/4thena-io/abyss/internal/api/rest/handler"
	"github.com/4thena-io/abyss/internal/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(
	app *handler.AppHandler,
	project *handler.ProjectHandler,
	template *handler.TemplateHandler,
	repo *handler.RepoHandler,
	team *handler.TeamHandler,
) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	web := web.NewFrontendHandler()

	r.Route("/", func(r chi.Router) {
		r.Get("/*", web.Serve)

		r.Route("/api", func(r chi.Router) {
			r.Route("/apps", func(r chi.Router) {
				r.Get("/", app.GetAllApps)
				r.Post("/", app.CreateApp)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.GetAppByID)
					r.Delete("/", app.DeleteApp)
					r.Get("/builds", app.GetAppBuilds)
					r.Get("/deployments", app.GetAppDeployments)
				})
			})

			r.Route("/projects", func(r chi.Router) {
				r.Get("/", project.GetAllProjects)
				r.Post("/", project.CreateProject)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", project.GetProjectByID)
					r.Delete("/", project.DeleteProject)
					r.Get("/apps", project.GetProjectApps)
				})
			})

			r.Route("/templates", func(r chi.Router) {
				r.Get("/", template.GetAllTemplates)
				r.Post("/", template.CreateTemplate)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", template.GetTemplateByName)
					r.Delete("/", template.DeleteTemplate)
				})
			})

			r.Route("/repos", func(r chi.Router) {
				r.Get("/", repo.GetAllRepos)
			})

			r.Route("/teams", func(r chi.Router) {
				r.Get("/", team.GetAllTeams)
				r.Post("/", team.CreateTeam)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", team.GetTeamByID)
					r.Delete("/", team.DeleteTeam)
					r.Get("/projects", team.GetTeamProjects)
					r.Get("/members", team.GetTeamMembers)
					r.Post("/members", team.AddTeamMember)
				})
			})
		})
	})

	return r
}
