package forgejo

import (
	"context"
	"fmt"
	"strings"

	fg "codeberg.org/mvdkleijn/forgejo-sdk/forgejo"
	"github.com/4thena-io/abyss/internal/model"
)

type ForgejoForge struct {
	client *fg.Client
	url    string
}

func NewForgejoForge(url, secret string) (*ForgejoForge, error) {
	client, err := fg.NewClient(url, fg.SetToken(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to create Forgejo client: %w", err)
	}
	return &ForgejoForge{client: client, url: url}, nil
}

func (f *ForgejoForge) OAuthEndpoints() (authPath, tokenPath string) {
	return "/login/oauth/authorize", "/login/oauth/access_token"
}

func (f *ForgejoForge) GetAuthenticatedUser(ctx context.Context) (*model.ForgeUser, error) {
	user, _, err := f.client.GetMyUserInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}
	name := user.FullName
	if name == "" {
		name = user.UserName
	}
	return &model.ForgeUser{
		ID:        user.ID,
		Username:  user.UserName,
		FullName:  name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}

func (f *ForgejoForge) GetUserByToken(ctx context.Context, token string) (*model.ForgeUser, error) {
	client, err := fg.NewClient(f.url, fg.SetToken(token))
	if err != nil {
		return nil, fmt.Errorf("failed to create user client: %w", err)
	}
	user, _, err := client.GetMyUserInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	name := user.FullName
	if name == "" {
		name = user.UserName
	}
	return &model.ForgeUser{
		ID:        user.ID,
		Username:  user.UserName,
		FullName:  name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}

func (f *ForgejoForge) IsMemberOfOwner(ctx context.Context, owner, username string) (bool, error) {
	member, _, err := f.client.CheckOrgMembership(owner, username)
	if err != nil {
		return false, fmt.Errorf("failed to check org membership: %w", err)
	}
	return member, nil
}

func (f *ForgejoForge) CreateRepo(ctx context.Context, owner, name string) (*model.Repo, error) {
	repo, _, err := f.client.CreateOrgRepo(owner, fg.CreateRepoOption{
		Name:    name,
		Private: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}
	return &model.Repo{
		ID:       repo.ID,
		Name:     repo.Name,
		FullName: repo.FullName,
		URL:      repo.HTMLURL,
		CloneURL: repo.CloneURL,
	}, nil
}

func (f *ForgejoForge) GetOrgRepos(ctx context.Context, name string) ([]model.Repo, error) {
	repos, _, err := f.client.ListOrgRepos(name, fg.ListOrgReposOptions{})
	if err != nil {
		return nil, err
	}

	res := make([]model.Repo, len(repos))
	for i, repo := range repos {
		res[i] = model.Repo{
			ID:       repo.ID,
			Name:     repo.Name,
			FullName: repo.FullName,
			URL:      repo.HTMLURL,
			CloneURL: repo.CloneURL,
		}
	}

	return res, nil
}

func (f *ForgejoForge) GetRepo(ctx context.Context, id int64) (*model.Repo, error) {
	repo, _, err := f.client.GetRepoByID(id)
	if err != nil {
		return nil, err
	}

	return &model.Repo{
		ID:       repo.ID,
		Name:     repo.Name,
		FullName: repo.FullName,
		URL:      repo.HTMLURL,
		CloneURL: repo.CloneURL,
	}, nil
}

func (f *ForgejoForge) DeleteRepo(ctx context.Context, owner, name string) error {
	_, err := f.client.DeleteRepo(owner, name)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}
	return nil
}

func (f *ForgejoForge) DeleteWebhook(ctx context.Context, owner, repo, callbackURL string) error {
	base, _, _ := strings.Cut(callbackURL, "?")
	hooks, _, err := f.client.ListRepoHooks(owner, repo, fg.ListHooksOptions{})
	if err != nil {
		return fmt.Errorf("failed to list webhooks: %w", err)
	}
	for _, hook := range hooks {
		u, _, _ := strings.Cut(hook.Config["url"], "?")
		if u == base {
			if _, err := f.client.DeleteRepoHook(owner, repo, hook.ID); err != nil {
				return fmt.Errorf("failed to delete webhook %d: %w", hook.ID, err)
			}
		}
	}
	return nil
}

func (f *ForgejoForge) CreateWebhook(ctx context.Context, owner, repo, callbackURL, secret, branch string) error {
	active := true
	_, _, err := f.client.CreateRepoHook(owner, repo, fg.CreateHookOption{
		Type: fg.HookTypeGitea,
		Config: map[string]string{
			"url":          callbackURL,
			"content_type": "json",
			"secret":       secret,
		},
		Events:       []string{"push"},
		BranchFilter: branch,
		Active:       active,
	})
	if err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}
	return nil
}
