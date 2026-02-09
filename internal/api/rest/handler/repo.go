package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/integration/forge"
	"github.com/4thena-io/abyss/internal/service"
)

type RepoHandler struct {
	service service.RepoService
}

func NewRepoHandler() *RepoHandler {
	forge, err := forge.NewForge()
	if err != nil {
		log.Fatalf("failed to create a git provider: %s", err)
	}
	service := service.NewRepoService(forge)
	return &RepoHandler{
		*service,
	}
}

func (h *RepoHandler) GetAllRepos(w http.ResponseWriter, r *http.Request) {
	repos, err := h.service.GetRepos(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	json.NewEncoder(w).Encode(res)
}
