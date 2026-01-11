package woodpecker

import (
	"context"
	"fmt"

	"git.4thena.io/4thena/abys/internal/dto/request"
	"git.4thena.io/4thena/abys/internal/dto/response"
	"go.woodpecker-ci.org/woodpecker/v3/woodpecker-go/woodpecker"
	"golang.org/x/oauth2"
)

type WoodpeckerCi struct {
	client woodpecker.Client
	host   string
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

func (c *WoodpeckerCi) ActivateRepo(ctx context.Context, req request.ActivateRepo) (*response.ActivateRepo, error) {
	repo, err := c.client.RepoPost(woodpecker.RepoPostOptions{
		ForgeRemoteID: req.ForgeRemoteId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to activate ci repo: %w", err)
	}
	return &response.ActivateRepo{
		RepoId:  repo.ID,
		RepoUrl: fmt.Sprintf("%s/repos/%d", c.host, repo.ID),
	}, nil
}

func (c *WoodpeckerCi) GetBuilds(ctx context.Context, repoID int64) ([]response.Build, error) {
	builds, err := c.client.PipelineList(repoID, woodpecker.PipelineListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch builds: %w", err)
	}

	res := make([]response.Build, len(builds))
	for i, build := range builds {
		res[i] = response.Build{
			ID:     build.ID,
			Number: build.Number,
			Status: build.Status,
			Branch: build.Branch,
			Commit: build.Commit,
			Link:   fmt.Sprintf("%s/repos/%d/pipeline/%d", c.host, repoID, build.Number),
		}
	}

	return res, nil
}
