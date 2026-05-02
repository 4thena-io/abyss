package forge

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
)

type Forge interface {
	GetAuthenticatedUser(ctx context.Context) (*model.ForgeUser, error)
	GetUserByToken(ctx context.Context, token string) (*model.ForgeUser, error)
	IsMemberOfOwner(ctx context.Context, owner, username string) (bool, error)
	CreateRepo(ctx context.Context, owner, name string) (*model.Repo, error)
	GetRepo(ctx context.Context, id int64) (*model.Repo, error)
	GetOrgRepos(ctx context.Context, name string) ([]model.Repo, error)
	DeleteRepo(ctx context.Context, owner, name string) error
}
