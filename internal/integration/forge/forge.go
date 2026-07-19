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
	CreateWebhook(ctx context.Context, owner, repo, callbackURL, secret, branch string) error
	// DeleteWebhook removes any webhooks whose URL matches callbackURL (query
	// params ignored). No-op if none found.
	DeleteWebhook(ctx context.Context, owner, repo, callbackURL string) error
	// OAuthEndpoints returns the forge's OAuth2 authorization and token URL
	// paths (relative to the forge host), e.g. "/login/oauth/authorize".
	OAuthEndpoints() (authPath, tokenPath string)
}
