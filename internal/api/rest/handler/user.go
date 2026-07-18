package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	teamService *service.TeamService
}

func NewUserHandler(teamService *service.TeamService) *UserHandler {
	return &UserHandler{teamService: teamService}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}

	stats, err := h.teamService.GetUserStats(r.Context(), uint(id))
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, "user not found")
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		response.InternalError(w)
		return
	}

	teams := make([]response.UserTeamMember, len(stats.Teams))
	for i, t := range stats.Teams {
		teams[i] = response.UserTeamMember{
			ID:           t.Team.ID,
			Name:         t.Team.Name,
			Role:         t.Role,
			MemberCount:  t.MemberCount,
			ProjectCount: t.ProjectCount,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response.UserProfile{
		ID:           stats.User.ID,
		Username:     stats.User.Username,
		Email:        stats.User.Email,
		AvatarURL:    stats.User.AvatarURL,
		IsAdmin:      stats.User.IsAdmin,
		TeamCount:    stats.TeamCount,
		ProjectCount: stats.ProjectCount,
		AppCount:     stats.AppCount,
		Teams:        teams,
	})
}
