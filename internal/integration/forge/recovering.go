package forge

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/4thena-io/abyss/internal/config"
	"github.com/4thena-io/abyss/internal/model"
)

var ErrUnavailable = errors.New("forge temporarily unavailable")

// Recovering defers network-dependent SDK initialization until a request needs
// it. Failed initialization is never cached; a later request can recover.
type Recovering struct {
	mu     sync.Mutex
	cfg    config.ForgeConfig
	client Forge
	create func(config.ForgeConfig) (Forge, error)
}

func NewRecovering(cfg config.ForgeConfig) (*Recovering, error) {
	switch cfg.Type {
	case "gitea", "forgejo", "github", "gitlab":
	default:
		return nil, fmt.Errorf("unknown forge type %s", cfg.Type)
	}
	return &Recovering{cfg: cfg, create: NewForge}, nil
}
func (f *Recovering) provider(ctx context.Context) (Forge, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.client != nil {
		return f.client, nil
	}
	client, err := f.create(f.cfg)
	if err != nil {
		return nil, ErrUnavailable
	}
	f.client = client
	return client, nil
}
func (f *Recovering) OAuthEndpoints() (string, string) {
	if f.cfg.Type == "gitlab" {
		return "/oauth/authorize", "/oauth/token"
	}
	return "/login/oauth/authorize", "/login/oauth/access_token"
}
func (f *Recovering) GetUserByToken(ctx context.Context, token string) (*model.ForgeUser, error) {
	// Identity lookup uses only the user's token, never bot initialization.
	cfg := f.cfg
	cfg.Token = token
	p, err := f.create(cfg)
	if err != nil {
		return nil, ErrUnavailable
	}
	return p.GetUserByToken(ctx, token)
}
func (f *Recovering) GetAuthenticatedUser(ctx context.Context) (*model.ForgeUser, error) {
	p, err := f.provider(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetAuthenticatedUser(ctx)
}
func (f *Recovering) IsMemberOfOwner(ctx context.Context, owner, username string) (bool, error) {
	p, err := f.provider(ctx)
	if err != nil {
		return false, err
	}
	return p.IsMemberOfOwner(ctx, owner, username)
}
func (f *Recovering) CreateRepo(ctx context.Context, owner, name string) (*model.Repo, error) {
	p, err := f.provider(ctx)
	if err != nil {
		return nil, err
	}
	return p.CreateRepo(ctx, owner, name)
}
func (f *Recovering) GetRepo(ctx context.Context, id int64) (*model.Repo, error) {
	p, err := f.provider(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetRepo(ctx, id)
}
func (f *Recovering) GetOrgRepos(ctx context.Context, owner string) ([]model.Repo, error) {
	p, err := f.provider(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetOrgRepos(ctx, owner)
}
func (f *Recovering) DeleteRepo(ctx context.Context, owner, name string) error {
	p, err := f.provider(ctx)
	if err != nil {
		return err
	}
	return p.DeleteRepo(ctx, owner, name)
}
func (f *Recovering) CreateWebhook(ctx context.Context, owner, repo, url, secret, branch string) error {
	p, err := f.provider(ctx)
	if err != nil {
		return err
	}
	return p.CreateWebhook(ctx, owner, repo, url, secret, branch)
}
func (f *Recovering) DeleteWebhook(ctx context.Context, owner, repo, url string) error {
	p, err := f.provider(ctx)
	if err != nil {
		return err
	}
	return p.DeleteWebhook(ctx, owner, repo, url)
}
