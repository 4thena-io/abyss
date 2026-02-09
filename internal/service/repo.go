package service

import (
	"context"

	"github.com/4thena-io/abyss/internal/integration/forge"
	"github.com/4thena-io/abyss/internal/model"
)

type RepoService struct {
	forge forge.Forge
}

func NewRepoService(forge forge.Forge) *RepoService {
	return &RepoService{
		forge,
	}
}

func (s *RepoService) GetRepos(ctx context.Context) ([]model.Repo, error) {
	repos, err := s.forge.GetOrgRepos(ctx, "4thena")
	if err != nil {
		return nil, err
	}
	return repos, nil
}
