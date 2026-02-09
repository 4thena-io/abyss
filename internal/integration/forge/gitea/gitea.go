package gitea

import (
	"context"
	"fmt"

	gt "code.gitea.io/sdk/gitea"
	"github.com/4thena-io/abyss/internal/model"
)

type GiteaForge struct {
	client *gt.Client
}

func NewGiteaForge(url, secret string) (*GiteaForge, error) {
	client, err := gt.NewClient(url, gt.SetToken(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gitea client: %w", err)
	}
	return &GiteaForge{client}, nil
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
