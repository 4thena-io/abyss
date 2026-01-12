package service

import (
	"context"

	"git.4thena.io/4thena/abys/internal/integration/forge"
	"git.4thena.io/4thena/abys/internal/model"
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
