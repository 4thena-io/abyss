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
	// CreateWebhook registers a push-event webhook on the given repository.
	// The webhook fires on every push to main/master; path filtering (docs/,
	// .abyss.yaml) is handled server-side by the abyss callback handler.
	CreateWebhook(ctx context.Context, owner, repo, callbackURL, secret string) error
}
