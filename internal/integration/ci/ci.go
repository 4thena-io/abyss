package ci

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
)

type CI interface {
	ActivateRepo(ctx context.Context, forgeRemoteID int64, slug string) (*model.CIRepo, error)
	DeleteRepo(ctx context.Context, ciRepoID int64, slug string) error
	GetBuilds(ctx context.Context, ciRepoID int64, slug string) ([]model.Build, error)
}
