package gitea

import (
	"context"
	"fmt"

	gt "code.gitea.io/sdk/gitea"
	"github.com/4thena-io/abyss/internal/model"
)

type GiteaForge struct {
	client *gt.Client
	url    string
}

func NewGiteaForge(url, secret string) (*GiteaForge, error) {
	client, err := gt.NewClient(url, gt.SetToken(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gitea client: %w", err)
	}
	return &GiteaForge{client: client, url: url}, nil
}

func (f *GiteaForge) GetAuthenticatedUser(ctx context.Context) (*model.ForgeUser, error) {
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

func (f *GiteaForge) GetUserByToken(ctx context.Context, token string) (*model.ForgeUser, error) {
	client, err := gt.NewClient(f.url, gt.SetToken(token))
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

func (f *GiteaForge) IsMemberOfOwner(ctx context.Context, owner, username string) (bool, error) {
	member, _, err := f.client.CheckOrgMembership(owner, username)
	if err != nil {
		return false, fmt.Errorf("failed to check org membership: %w", err)
	}
	return member, nil
}

func (f *GiteaForge) CreateRepo(ctx context.Context, owner, name string) (*model.Repo, error) {
	repo, _, err := f.client.CreateOrgRepo(owner, gt.CreateRepoOption{
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

func (f *GiteaForge) GetOrgRepos(ctx context.Context, name string) ([]model.Repo, error) {
	repos, _, err := f.client.ListOrgRepos(name, gt.ListOrgReposOptions{})
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

func (f *GiteaForge) GetRepo(ctx context.Context, id int64) (*model.Repo, error) {
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

func (f *GiteaForge) DeleteRepo(ctx context.Context, owner, name string) error {
	_, err := f.client.DeleteRepo(owner, name)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}
	return nil
}
