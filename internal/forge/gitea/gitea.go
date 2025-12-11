package gitea

import (
	"context"
	"fmt"

	gt "code.gitea.io/sdk/gitea"
	"git.4thena.io/4thena/abys/internal/dto"
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

func (f *GiteaForge) CreateRepo(ctx context.Context, owner, name string) (*dto.CreateRepoResponseDTO, error) {
	repo, _, err := f.client.CreateOrgRepo(owner, gt.CreateRepoOption{
		Name: name,
		Private: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}
	return &dto.CreateRepoResponseDTO{
		RepoId: repo.ID,
		CloneUrl: repo.CloneURL,
		HtmlUrl: repo.HTMLURL,
	}, nil
}

func (f *GiteaForge) DeleteRepo(ctx context.Context, owner, name string) (error) {
	_, err := f.client.DeleteRepo(owner, name)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}
	return nil
}
