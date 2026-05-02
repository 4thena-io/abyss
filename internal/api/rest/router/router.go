package router

import (
	"net/http"

	"github.com/4thena-io/abyss/internal/api/rest/handler"
	"github.com/4thena-io/abyss/internal/api/rest/middleware"
	"github.com/4thena-io/abyss/internal/web"
	"github.com/go-chi/chi/v5"
)

func New(
	requireSetup func(http.Handler) http.Handler,
	requireAuth func(http.Handler) http.Handler,
	setupHandler *handler.SetupHandler,
	auth *handler.AuthHandler,
	app *handler.AppHandler,
	project *handler.ProjectHandler,
	template *handler.TemplateHandler,
	repo *handler.RepoHandler,
	team *handler.TeamHandler,
) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	frontend := web.NewFrontendHandler()

	// Public setup endpoints — always accessible, no auth required
	r.Route("/api/setup", func(r chi.Router) {
		r.Get("/status", setupHandler.Status)
		r.Post("/configure", setupHandler.Configure)
	})

	// Auth routes: require OAuth to be configured
	r.Route("/auth", func(r chi.Router) {
		r.Use(requireSetup)
		r.Get("/login", auth.Login)
		r.Get("/callback", auth.Callback)
		r.Post("/logout", auth.Logout)
	})

	// Protected API routes: require setup + authentication
	r.Route("/api", func(r chi.Router) {
		r.Use(requireSetup)
		r.Use(requireAuth)

		r.Get("/me", auth.Me)

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

	// Frontend SPA — catch-all, always accessible
	r.Get("/*", frontend.Serve)

	return r
}
