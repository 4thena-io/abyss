package git

import (
	"context"
	"fmt"

	gt "code.gitea.io/sdk/gitea"
)

type GiteaProvider struct {
	client *gt.Client
}

func NewGiteaProvider(url, secret string) (Provider, error) {
	client, err := gt.NewClient(url, gt.SetToken(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gitea client: %w", err)
	}
	return &GiteaProvider{client}, nil
}

func (g *GiteaProvider) CreateRepo(ctx context.Context, orgName, name string) (*CreateRepoResponseDTO, error) {
	repo, _, err := g.client.CreateOrgRepo(orgName, gt.CreateRepoOption{
		Name: name,
		Private: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}
	return &CreateRepoResponseDTO{
		RepoId: repo.ID,
		CloneUrl: repo.CloneURL,
		HtmlUrl: repo.HTMLURL,
	}, nil
}
