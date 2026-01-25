package router

import (
	"git.4thena.io/4thena/abys/internal/handler"
	"git.4thena.io/4thena/abys/internal/web"
	"github.com/go-chi/chi/v5"
)

func New() chi.Router {
	r := chi.NewRouter()

	web := web.NewFrontendHandler()
	app := handler.NewAppHandler()
	project := handler.NewProjectHandler()
	template := handler.NewTemplateHandler()
	repo := handler.NewRepoHandler()

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
		})
	})

	return r
}
