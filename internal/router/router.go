package router

import (
	"git.4thena.io/4thena/abys/internal/handler"
	"git.4thena.io/4thena/abys/internal/web"
	"github.com/gorilla/mux"
)

func New() *mux.Router {
	router := mux.NewRouter()

	// API
	api := router.PathPrefix("/api").Subrouter()
	handler.NewAppHandler().RegisterRoutes(api.PathPrefix("/apps").Subrouter())
	handler.NewProjectHandler().RegisterRoutes(api.PathPrefix("/projects").Subrouter())
	handler.NewTemplateHandler().RegisterRoutes(api.PathPrefix("/templates").Subrouter())

	// Frontend
	web.RegisterRoutes(router)

	return router
}
