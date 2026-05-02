package github

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
)

// TODO: Implement GitHub forge when GitHub support is prioritized

type GithubForge struct {
}

func NewGithubForge(url, secret string) (*GithubForge, error) {
	return &GithubForge{}, nil
}

func (f *GithubForge) GetAuthenticatedUser(ctx context.Context) (*model.ForgeUser, error) {
	// TODO: implement when GitHub support is prioritized
	return &model.ForgeUser{
		Username: "github-bot",
		FullName: "GitHub Bot",
		Email:    "bot@github.com",
	}, nil
}

func (f *GithubForge) GetUserByToken(ctx context.Context, token string) (*model.ForgeUser, error) {
	// TODO: implement when GitHub support is prioritized
	return &model.ForgeUser{}, nil
}

func (f *GithubForge) IsMemberOfOwner(ctx context.Context, owner, username string) (bool, error) {
	// TODO: implement when GitHub support is prioritized
	return true, nil
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
