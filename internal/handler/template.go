package handler

import (
	"encoding/json"
	"net/http"

	"git.4thena.io/4thena/abys/internal/database"
	"git.4thena.io/4thena/abys/internal/dto/request"
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

	app, err := h.service.SaveTemplate(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app)
}

func (h *TemplateHandler) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	apps, err := h.service.GetAllTemplates(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
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
	json.NewEncoder(w).Encode(template)
}

func (h *TemplateHandler) RegisterTemplateRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/templates").Subrouter()

	userRouter.HandleFunc("", h.CreateTemplate).Methods("POST")
	userRouter.HandleFunc("", h.GetAllTemplates).Methods("GET")
	userRouter.HandleFunc("/{value}", h.GetTemplateByName).Methods("GET")
}
