package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/4thena-io/abyss/internal/api/rest/request"
	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/constant"
	"github.com/4thena-io/abyss/internal/database"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/repository"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/go-chi/chi/v5"
)

type TeamHandler struct {
	service service.TeamService
}

func NewTeamHandler() *TeamHandler {
	teamRepository := repository.NewTeamRepository(database.Connection)
	projectRepository := repository.NewProjectRepository(database.Connection)
	appRepository := repository.NewAppRepository(database.Connection)

	teamService := service.NewTeamService(*teamRepository, *projectRepository, *appRepository)
	return &TeamHandler{*teamService}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req request.CreateTeam
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	team, err := h.service.SaveTeam(r.Context(), &model.Team{
		Name:        req.Name,
		Description: req.Description,
		Status:      constant.StatusActive,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.Team{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
	})
}

func (h *TeamHandler) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.service.GetAllTeams(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := make([]response.Team, len(teams))
	for i, t := range teams {
		res[i] = response.Team{
			ID:           t.Team.ID,
			Name:         t.Team.Name,
			Description:  t.Team.Description,
			MemberCount:  t.MemberCount,
			ProjectCount: t.ProjectCount,
			AppCount:     t.AppCount,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *TeamHandler) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	t, err := h.service.GetTeamByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response.Team{
		ID:           t.Team.ID,
		Name:         t.Team.Name,
		Description:  t.Team.Description,
		MemberCount:  t.MemberCount,
		ProjectCount: t.ProjectCount,
		AppCount:     t.AppCount,
	})
}

func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTeam(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) GetTeamProjects(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	projects, err := h.service.GetTeamProjects(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := make([]response.Project, len(projects))
	for i, p := range projects {
		res[i] = response.Project{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			TeamID:      p.TeamID,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *TeamHandler) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	members, err := h.service.GetTeamMembers(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := make([]response.TeamMember, len(members))
	for i, m := range members {
		res[i] = response.TeamMember{
			ID:       m.ID,
			TeamID:   m.TeamID,
			Username: m.Username,
			Role:     m.Role,
			JoinedAt: m.CreatedAt.Format(time.RFC3339),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *TeamHandler) AddTeamMember(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req request.AddTeamMember
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	role := req.Role
	if role == "" {
		role = "member"
	}

	member, err := h.service.AddTeamMember(r.Context(), &model.TeamMember{
		TeamID:   uint(id),
		Username: req.Username,
		Role:     role,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.TeamMember{
		ID:       member.ID,
		TeamID:   member.TeamID,
		Username: member.Username,
		Role:     member.Role,
		JoinedAt: member.CreatedAt.Format(time.RFC3339),
	})
}
