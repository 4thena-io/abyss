package github

import (
	"context"

	"git.4thena.io/4thena/abys/internal/model"
)

// TODO: Implement GitHub forge when GitHub support is prioritized

type GithubForge struct {
}

func NewGithubForge(url, secret string) (*GithubForge, error) {
	return &GithubForge{}, nil
}

func (f *GithubForge) CreateRepo(ctx context.Context, owner, name string) (*model.Repo, error) {
	return &model.Repo{
		ID:       1,
		URL:      "https://github.com/owner/repo",
		CloneURL: "https://github.com/owner/repo.git",
	}, nil
}

func (f *GithubForge) GetRepo(ctx context.Context, id int64) (*model.Repo, error) {
	return nil, nil
}

func (f *GithubForge) GetOrgRepos(ctx context.Context, name string) ([]model.Repo, error) {
	return nil, nil
}

func (f *GithubForge) DeleteRepo(ctx context.Context, owner, name string) error {
	return nil
}
