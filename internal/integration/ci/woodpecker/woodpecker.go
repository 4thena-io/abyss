package woodpecker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/4thena-io/abyss/internal/model"
	"go.woodpecker-ci.org/woodpecker/v3/woodpecker-go/woodpecker"
	"golang.org/x/oauth2"
)

type WoodpeckerCi struct {
	client     woodpecker.Client
	httpClient *http.Client
	host       string
}

func NewWoodpeckerCi(url, token string) (*WoodpeckerCi, error) {
	cfg := new(oauth2.Config)
	auth := cfg.Client(
		context.Background(),
		&oauth2.Token{
			AccessToken: token,
		},
	)
	url = strings.TrimRight(url, "/")
	client := woodpecker.NewClient(url, auth)
	return &WoodpeckerCi{
		client:     client,
		httpClient: auth,
		host:       url,
	}, nil
}

// ActivateRepo activates CI for a repository using the forge remote ID.
// The slug parameter is ignored by Woodpecker.
func (c *WoodpeckerCi) ActivateRepo(ctx context.Context, forgeRemoteID int64, slug string) (*model.CIRepo, error) {
	repo, err := c.client.RepoPost(woodpecker.RepoPostOptions{
		ForgeRemoteID: forgeRemoteID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to activate ci repo: %w", err)
	}

	return &model.CIRepo{
		ID:   repo.ID,
		URL:  fmt.Sprintf("%s/repos/%d", c.host, repo.ID),
		Slug: repo.FullName,
	}, nil
}

// DeleteRepo permanently removes a repository and its CI data using the CI repo ID.
// The slug parameter is ignored by Woodpecker.
func (c *WoodpeckerCi) DeleteRepo(ctx context.Context, ciRepoID int64, slug string) error {
	// The SDK's RepoDel omits remove=true and only deactivates the repository.
	url := fmt.Sprintf("%s/api/repos/%d?remove=true", c.host, ciRepoID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create ci repo deletion request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete ci repo: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("failed to delete ci repo: %s", resp.Status)
	}
	return nil
}

// GetBuilds returns pipelines for a repository using the CI repo ID.
// The slug parameter is ignored by Woodpecker.
func (c *WoodpeckerCi) GetBuilds(ctx context.Context, ciRepoID int64, slug string) ([]model.Build, error) {
	builds, err := c.client.PipelineList(ciRepoID, woodpecker.PipelineListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch builds: %w", err)
	}

	res := make([]model.Build, len(builds))
	for i, build := range builds {
		res[i] = model.Build{
			ID:     build.ID,
			Number: build.Number,
			Status: build.Status,
			Branch: build.Branch,
			Commit: build.Commit,
			Link:   fmt.Sprintf("%s/repos/%d/pipeline/%d", c.host, ciRepoID, build.Number),
		}
	}

	return res, nil
}
