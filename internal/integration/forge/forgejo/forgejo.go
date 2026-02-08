package forgejo

import (
	"context"
	"fmt"

	fg "codeberg.org/mvdkleijn/forgejo-sdk/forgejo"
	"git.4thena.io/4thena/abys/internal/model"
)

type ForgejoForge struct {
	client *fg.Client
}

func NewForgejoForge(url, secret string) (*ForgejoForge, error) {
	client, err := fg.NewClient(url, fg.SetToken(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to create Forgejo client: %w", err)
	}
	return &ForgejoForge{client}, nil
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
