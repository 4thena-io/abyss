package gitea_actions

import (
	"context"
	"fmt"
	"strings"

	gt "code.gitea.io/sdk/gitea"
	"github.com/4thena-io/abyss/internal/model"
)

type GiteaActionsCI struct {
	client *gt.Client
	host string
}

// parseSlug splits "owner/repo" into owner and repo parts.
func parseSlug(slug string) (owner, repo string, err error) {
	parts := strings.SplitN(slug, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid slug format: %s", slug)
	}
	return parts[0], parts[1], nil
}

func NewGiteaActionsCI(url, token string) (*GiteaActionsCI, error) {
	client, err := gt.NewClient(url, gt.SetToken(token))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gitea client: %w", err)
	}
	return &GiteaActionsCI{
		client: client,
		host: url,
	}, nil
}

func (c *GiteaActionsCI) ActivateRepo(ctx context.Context, forgeRemoteID int64, slug string) (*model.CIRepo, error) {
	return &model.CIRepo{
		ID: forgeRemoteID,
		Slug: slug,
		URL: fmt.Sprintf("%s/%s/%s", c.host, slug, "actions"),
	}, nil
}

func (c *GiteaActionsCI) DeleteRepo(ctx context.Context, ciRepoID int64, slug string) error {
	return nil
}

func (c *GiteaActionsCI) GetBuilds(ctx context.Context, ciRepoID int64, slug string) ([]model.Build, error) {
	owner, repo, err := parseSlug(slug)
	if err != nil {
		return nil, err
	}

	runs, _, err := c.client.ListRepoActionRuns(
		owner,
		repo,
		gt.ListRepoActionRunsOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch builds: %w", err)
	}

	res := make([]model.Build, len(runs.WorkflowRuns))

	for i, run := range runs.WorkflowRuns {
		res[i] = model.Build{
			ID:     run.ID,
			Number: int64(run.RunNumber),
			Status: run.Status,
			Branch: run.HeadBranch,
			Commit: run.HeadSha,
			Link:   run.HTMLURL,
		}
	}

	return res, nil
}
