package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	gh "github.com/google/go-github/v67/github"
	"github.com/4thena-io/abyss/internal/model"
	"golang.org/x/oauth2"
)

type GithubForge struct {
	client *gh.Client
	host   string
}

func NewGithubForge(host, token string) (*GithubForge, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)

	var client *gh.Client
	if host != "" && host != "https://github.com" && host != "https://api.github.com" {
		// GitHub Enterprise
		var err error
		client, err = gh.NewClient(tc).WithEnterpriseURLs(host+"/api/v3/", host+"/api/uploads/")
		if err != nil {
			return nil, fmt.Errorf("failed to create GitHub Enterprise client: %w", err)
		}
	} else {
		client = gh.NewClient(tc)
	}

	return &GithubForge{client: client, host: host}, nil
}

func (f *GithubForge) OAuthEndpoints() (authPath, tokenPath string) {
	return "/login/oauth/authorize", "/login/oauth/access_token"
}

func (f *GithubForge) GetAuthenticatedUser(ctx context.Context) (*model.ForgeUser, error) {
	user, _, err := f.client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}
	name := user.GetName()
	if name == "" {
		name = user.GetLogin()
	}
	return &model.ForgeUser{
		ID:        user.GetID(),
		Username:  user.GetLogin(),
		FullName:  name,
		Email:     user.GetEmail(),
		AvatarURL: user.GetAvatarURL(),
	}, nil
}

func (f *GithubForge) GetUserByToken(ctx context.Context, token string) (*model.ForgeUser, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	var client *gh.Client
	if f.host != "" && f.host != "https://github.com" && f.host != "https://api.github.com" {
		var err error
		client, err = gh.NewClient(tc).WithEnterpriseURLs(f.host+"/api/v3/", f.host+"/api/uploads/")
		if err != nil {
			return nil, fmt.Errorf("failed to create user client: %w", err)
		}
	} else {
		client = gh.NewClient(tc)
	}

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	name := user.GetName()
	if name == "" {
		name = user.GetLogin()
	}
	return &model.ForgeUser{
		ID:        user.GetID(),
		Username:  user.GetLogin(),
		FullName:  name,
		Email:     user.GetEmail(),
		AvatarURL: user.GetAvatarURL(),
	}, nil
}

func (f *GithubForge) IsMemberOfOwner(ctx context.Context, owner, username string) (bool, error) {
	_, resp, err := f.client.Organizations.IsMember(ctx, owner, username)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to check org membership: %w", err)
	}
	return true, nil
}

func (f *GithubForge) CreateRepo(ctx context.Context, owner, name string) (*model.Repo, error) {
	private := false
	repo, _, err := f.client.Repositories.Create(ctx, owner, &gh.Repository{
		Name:    gh.String(name),
		Private: &private,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}
	return &model.Repo{
		ID:       repo.GetID(),
		Name:     repo.GetName(),
		FullName: repo.GetFullName(),
		URL:      repo.GetHTMLURL(),
		CloneURL: repo.GetCloneURL(),
	}, nil
}

func (f *GithubForge) GetRepo(ctx context.Context, id int64) (*model.Repo, error) {
	// GitHub API requires owner+repo, not ID for most endpoints.
	// We use the search API to find by ID.
	repo, _, err := f.client.Repositories.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository by ID: %w", err)
	}
	return &model.Repo{
		ID:       repo.GetID(),
		Name:     repo.GetName(),
		FullName: repo.GetFullName(),
		URL:      repo.GetHTMLURL(),
		CloneURL: repo.GetCloneURL(),
	}, nil
}

func (f *GithubForge) GetOrgRepos(ctx context.Context, name string) ([]model.Repo, error) {
	var allRepos []*gh.Repository
	opts := &gh.RepositoryListByOrgOptions{ListOptions: gh.ListOptions{PerPage: 100}}
	for {
		repos, resp, err := f.client.Repositories.ListByOrg(ctx, name, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list org repos: %w", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	res := make([]model.Repo, len(allRepos))
	for i, repo := range allRepos {
		res[i] = model.Repo{
			ID:       repo.GetID(),
			Name:     repo.GetName(),
			FullName: repo.GetFullName(),
			URL:      repo.GetHTMLURL(),
			CloneURL: repo.GetCloneURL(),
		}
	}
	return res, nil
}

func (f *GithubForge) DeleteRepo(ctx context.Context, owner, name string) error {
	_, err := f.client.Repositories.Delete(ctx, owner, name)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}
	return nil
}

func (f *GithubForge) DeleteWebhook(ctx context.Context, owner, repo, callbackURL string) error {
	base, _, _ := strings.Cut(callbackURL, "?")
	hooks, _, err := f.client.Repositories.ListHooks(ctx, owner, repo, nil)
	if err != nil {
		return fmt.Errorf("failed to list webhooks: %w", err)
	}
	for _, hook := range hooks {
		u, _, _ := strings.Cut(hook.GetConfig().GetURL(), "?")
		if u == base {
			if _, err := f.client.Repositories.DeleteHook(ctx, owner, repo, hook.GetID()); err != nil {
				return fmt.Errorf("failed to delete webhook %d: %w", hook.GetID(), err)
			}
		}
	}
	return nil
}

func (f *GithubForge) CreateWebhook(ctx context.Context, owner, repo, callbackURL, secret, branch string) error {
	active := true
	_, _, err := f.client.Repositories.CreateHook(ctx, owner, repo, &gh.Hook{
		Config: &gh.HookConfig{
			URL:         gh.String(callbackURL),
			ContentType: gh.String("json"),
			Secret:      gh.String(secret),
		},
		Events: []string{"push"},
		Active: &active,
	})
	if err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}
	return nil
}
