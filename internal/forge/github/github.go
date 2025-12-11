package github

import (
	"context"

	"git.4thena.io/4thena/abys/internal/dto"
)

type GithubForge struct {
}

func NewGithubForge(url, secret string) (*GithubForge, error) {
	return &GithubForge{}, nil
}

func (f *GithubForge) CreateRepo(ctx context.Context, owner, name string) (*dto.CreateRepoResponseDTO, error) {
	return &dto.CreateRepoResponseDTO{
		RepoId: 1,
		CloneUrl: "https://github.com/owner/repo.git",
		HtmlUrl: "https://github.com/owner/repo",
	}, nil
}

func (f *GithubForge) DeleteRepo(ctx context.Context, owner, name string) (error) {
	return nil
}
