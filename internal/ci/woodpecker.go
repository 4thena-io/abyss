package ci

import (
	"context"
	"fmt"

	"git.d4ramirez.com/project-abyss/abys-api/internal/dto"
	"go.woodpecker-ci.org/woodpecker/v3/woodpecker-go/woodpecker"
	"golang.org/x/oauth2"
)


type WoodpeckerCi struct{
	client	woodpecker.Client
	host		string
}

func NewWoodpeckerCi(url, token string) (*WoodpeckerCi, error) {
	cfg := new(oauth2.Config)
	auth := cfg.Client(
		context.Background(),
		&oauth2.Token{
			AccessToken: token,
		},
	)
	client := woodpecker.NewClient(url, auth)
	return &WoodpeckerCi{
		client,
		url,
	}, nil 
}

func (c *WoodpeckerCi) ActivateRepo(ctx context.Context, req dto.ActivateRepoRequestDTO) (*dto.ActivateRepoResponseDTO, error) {
	repo, err := c.client.RepoPost(woodpecker.RepoPostOptions{
		ForgeRemoteID: req.ForgeRemoteId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to activate ci repo: %w", err)
	}
	return &dto.ActivateRepoResponseDTO{
		RepoId: repo.ID,
		RepoUrl: fmt.Sprintf("%s/repos/%d", c.host, repo.ID),
	}, nil
}
