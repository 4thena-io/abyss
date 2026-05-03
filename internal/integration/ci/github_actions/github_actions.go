package github_actions

import (
	"context"
	"fmt"
	"strings"

	gh "github.com/google/go-github/v67/github"
	"github.com/4thena-io/abyss/internal/model"
	"golang.org/x/oauth2"
)

type GithubActionsCI struct {
	client *gh.Client
	host   string
}

func parseSlug(slug string) (owner, repo string, err error) {
	parts := strings.SplitN(slug, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid slug format: %s", slug)
	}
	return parts[0], parts[1], nil
}

func NewGithubActionsCI(host, token string) (*GithubActionsCI, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)

	var client *gh.Client
	if host != "" && host != "https://github.com" && host != "https://api.github.com" {
		var err error
		client, err = gh.NewClient(tc).WithEnterpriseURLs(host+"/api/v3/", host+"/api/uploads/")
		if err != nil {
			return nil, fmt.Errorf("failed to create GitHub Enterprise client: %w", err)
		}
	} else {
		client = gh.NewClient(tc)
	}

	return &GithubActionsCI{client: client, host: host}, nil
}

func (c *GithubActionsCI) ActivateRepo(ctx context.Context, forgeRemoteID int64, slug string) (*model.CIRepo, error) {
	// GitHub Actions is always available; enable Actions on the repo if needed.
	owner, repo, err := parseSlug(slug)
	if err != nil {
		return nil, err
	}
	return &model.CIRepo{
		ID:   forgeRemoteID,
		Slug: slug,
		URL:  fmt.Sprintf("https://github.com/%s/%s/actions", owner, repo),
	}, nil
}

func (c *GithubActionsCI) DeleteRepo(ctx context.Context, ciRepoID int64, slug string) error {
	return nil
}

func (c *GithubActionsCI) GetBuilds(ctx context.Context, ciRepoID int64, slug string) ([]model.Build, error) {
	owner, repo, err := parseSlug(slug)
	if err != nil {
		return nil, err
	}

	runs, _, err := c.client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, &gh.ListWorkflowRunsOptions{
		ListOptions: gh.ListOptions{PerPage: 50},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workflow runs: %w", err)
	}

	res := make([]model.Build, len(runs.WorkflowRuns))
	for i, run := range runs.WorkflowRuns {
		res[i] = model.Build{
			ID:     run.GetID(),
			Number: int64(run.GetRunNumber()),
			Status: run.GetStatus(),
			Branch: run.GetHeadBranch(),
			Commit: run.GetHeadSHA(),
			Link:   run.GetHTMLURL(),
		}
	}
	return res, nil
}
