package drone

import (
	"context"
	"fmt"
	"strings"

	"git.4thena.io/4thena/abys/internal/model"
	"github.com/drone/drone-go/drone"
	"golang.org/x/oauth2"
)

type DroneCi struct {
	client drone.Client
	host   string
}

func NewDroneCi(url, token string) (*DroneCi, error) {
	cfg := new(oauth2.Config)
	auth := cfg.Client(
		context.Background(),
		&oauth2.Token{
			AccessToken: token,
		},
	)
	client := drone.NewClient(url, auth)
	return &DroneCi{
		client: client,
		host:   url,
	}, nil
}

// parseSlug splits "owner/repo" into owner and repo parts.
func parseSlug(slug string) (owner, repo string, err error) {
	parts := strings.SplitN(slug, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid slug format: %s", slug)
	}
	return parts[0], parts[1], nil
}

// ActivateRepo activates CI for a repository using the slug (owner/repo).
// The forgeRemoteID parameter is ignored by Drone.
func (c *DroneCi) ActivateRepo(ctx context.Context, forgeRemoteID int64, slug string) (*model.CIRepo, error) {
	owner, repo, err := parseSlug(slug)
	if err != nil {
		return nil, err
	}

	droneRepo, err := c.client.RepoEnable(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to activate ci repo: %w", err)
	}

	return &model.CIRepo{
		ID:   droneRepo.ID,
		Slug: droneRepo.Slug,
		URL:  fmt.Sprintf("%s/%s", c.host, droneRepo.Slug),
	}, nil
}

// DeleteRepo removes a repository from CI using the slug (owner/repo).
// The ciRepoID parameter is ignored by Drone.
func (c *DroneCi) DeleteRepo(ctx context.Context, ciRepoID int64, slug string) error {
	owner, repo, err := parseSlug(slug)
	if err != nil {
		return err
	}

	return c.client.RepoDisable(owner, repo)
}

// GetBuilds returns builds for a repository using the slug (owner/repo).
// The ciRepoID parameter is ignored by Drone.
func (c *DroneCi) GetBuilds(ctx context.Context, ciRepoID int64, slug string) ([]model.Build, error) {
	owner, repo, err := parseSlug(slug)
	if err != nil {
		return nil, err
	}

	builds, err := c.client.BuildList(owner, repo, drone.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch builds: %w", err)
	}

	res := make([]model.Build, len(builds))
	for i, build := range builds {
		res[i] = model.Build{
			ID:     build.ID,
			Number: build.Number,
			Status: build.Status,
			Branch: build.Source,
			Commit: build.After,
			Link:   fmt.Sprintf("%s/%s/%d", c.host, slug, build.Number),
		}
	}

	return res, nil
}
