package handler

import (
	"encoding/json"
	"net/http"

	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/rs/zerolog/log"
)

type RepoHandler struct {
	service service.RepoService
}

func NewRepoHandler(service *service.RepoService) *RepoHandler {
	return &RepoHandler{*service}
}

func (h *RepoHandler) GetAllRepos(w http.ResponseWriter, r *http.Request) {
	repos, err := h.service.GetRepos(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get repos")
		response.InternalError(w)
		return
	}

	res := make([]response.Repo, len(repos))
	for i, repo := range repos {
		res[i] = response.Repo{
			ID:       repo.ID,
			Name:     repo.Name,
			FullName: repo.FullName,
			URL:      repo.URL,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}
