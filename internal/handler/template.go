package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.4thena.io/4thena/abys/internal/database"
	"git.4thena.io/4thena/abys/internal/dto/request"
	"git.4thena.io/4thena/abys/internal/dto/response"
	"git.4thena.io/4thena/abys/internal/model"
	"git.4thena.io/4thena/abys/internal/repository"
	"git.4thena.io/4thena/abys/internal/service"
	"github.com/gorilla/mux"
)

type TemplateHandler struct {
	service service.TemplateService
}

func NewTemplateHandler() *TemplateHandler {
	repository := repository.NewTemplateRepository(database.Connection)
	service := service.NewTemplateService(*repository)
	return &TemplateHandler{*service}
}

func (h *TemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req request.CreateTemplate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	template, err := h.service.SaveTemplate(r.Context(), &model.Template{
		Name:        req.Name,
		Description: req.Description,
		Kind:        req.Kind,
		Language:    req.Language,
		RepoURL:     req.RepoURL,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.Template{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		Kind:        template.Kind,
		Language:    template.Language,
		RepoURL:     template.RepoURL,
	})
}

func (h *TemplateHandler) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.service.GetAllTemplates(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := make([]response.Template, len(templates))
	for i, template := range templates {
		res[i] = response.Template{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			Kind:        template.Kind,
			Language:    template.Language,
			RepoURL:     template.RepoURL,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *TemplateHandler) GetTemplateByName(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	template, err := h.service.GetTemplateByName(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if template == nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.Template{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		Kind:        template.Kind,
		Language:    template.Language,
		RepoURL:     template.RepoURL,
	})
}

func (h *TemplateHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	err = h.service.DeleteTemplate(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (h *TemplateHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("", h.CreateTemplate).Methods("POST")
	router.HandleFunc("", h.GetAllTemplates).Methods("GET")
	router.HandleFunc("/{value}", h.GetTemplateByName).Methods("GET")
	router.HandleFunc("/{id}", h.DeleteTemplate).Methods("DELETE")
}
