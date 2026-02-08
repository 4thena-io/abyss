package forge

import (
	"context"

	"git.4thena.io/4thena/abys/internal/model"
)

type Forge interface {
	CreateRepo(ctx context.Context, owner, name string) (*model.Repo, error)
	GetRepo(ctx context.Context, id int64) (*model.Repo, error)
	GetOrgRepos(ctx context.Context, name string) ([]model.Repo, error)
	DeleteRepo(ctx context.Context, owner, name string) error
}
