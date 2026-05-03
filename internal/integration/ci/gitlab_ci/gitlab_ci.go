package gitlab_ci

import (
	"context"
	"fmt"
	"strings"

	gl "gitlab.com/gitlab-org/api/client-go"
	"github.com/4thena-io/abyss/internal/model"
)

type GitlabCI struct {
	client *gl.Client
	host   string
}

func parseSlug(slug string) (string, error) {
	// GitLab uses "owner/repo" as the project path directly
	if !strings.Contains(slug, "/") {
		return "", fmt.Errorf("invalid slug format: %s", slug)
	}
	return slug, nil
}

func NewGitlabCI(host, token string) (*GitlabCI, error) {
	opts := []gl.ClientOptionFunc{gl.WithBaseURL(host)}
	client, err := gl.NewClient(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}
	return &GitlabCI{client: client, host: host}, nil
}

func (c *GitlabCI) ActivateRepo(ctx context.Context, forgeRemoteID int64, slug string) (*model.CIRepo, error) {
	// GitLab CI is built-in; no activation needed.
	return &model.CIRepo{
		ID:   forgeRemoteID,
		Slug: slug,
		URL:  fmt.Sprintf("%s/%s/-/pipelines", c.host, slug),
	}, nil
}

func (c *GitlabCI) DeleteRepo(ctx context.Context, ciRepoID int64, slug string) error {
	return nil
}

func (c *GitlabCI) GetBuilds(ctx context.Context, ciRepoID int64, slug string) ([]model.Build, error) {
	pid, err := parseSlug(slug)
	if err != nil {
		return nil, err
	}

	pipelines, _, err := c.client.Pipelines.ListProjectPipelines(pid, &gl.ListProjectPipelinesOptions{
		ListOptions: gl.ListOptions{PerPage: 50},
	}, gl.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pipelines: %w", err)
	}

	res := make([]model.Build, len(pipelines))
	for i, p := range pipelines {
		res[i] = model.Build{
			ID:     int64(p.ID),
			Number: int64(p.IID),
			Status: p.Status,
			Branch: p.Ref,
			Commit: p.SHA,
			Link:   p.WebURL,
		}
	}
	return res, nil
}
