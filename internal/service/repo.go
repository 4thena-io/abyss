package service

import (
	"context"

	"github.com/4thena-io/abyss/internal/integration/forge"
	"github.com/4thena-io/abyss/internal/model"
)

type RepoService struct {
	forge forge.Forge
	owner string
}

func NewRepoService(forge forge.Forge, owner string) *RepoService {
	return &RepoService{forge, owner}
}

func (s *RepoService) GetRepos(ctx context.Context) ([]model.Repo, error) {
	return s.forge.GetOrgRepos(ctx, s.owner)
}
