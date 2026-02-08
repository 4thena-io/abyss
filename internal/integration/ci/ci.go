package ci

import (
	"context"

	"git.4thena.io/4thena/abys/internal/model"
)

type CI interface {
	ActivateRepo(ctx context.Context, forgeRemoteID int64, slug string) (*model.CIRepo, error)
	DeleteRepo(ctx context.Context, ciRepoID int64, slug string) error
	GetBuilds(ctx context.Context, ciRepoID int64, slug string) ([]model.Build, error)
}
