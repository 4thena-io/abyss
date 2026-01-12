package jenkins

import (
	"context"

	"git.4thena.io/4thena/abys/internal/dto/request"
	"git.4thena.io/4thena/abys/internal/dto/response"
)

type JenkinsCi struct{}

func NewJenkinsCi(url, token string) (*JenkinsCi, error) {
	return &JenkinsCi{}, nil
}

func (c *JenkinsCi) ActivateRepo(ctx context.Context, req request.ActivateRepo) (*response.ActivateRepo, error) {
	return &response.ActivateRepo{
		RepoId:  1,
		RepoUrl: "https://jenkins.com/1",
	}, nil
}

func (c *JenkinsCi) GetBuilds(ctx context.Context, repoId int64) ([]response.Build, error) {
	res := make([]response.Build, 1)
	res[0] = response.Build{}
	return res, nil
}
