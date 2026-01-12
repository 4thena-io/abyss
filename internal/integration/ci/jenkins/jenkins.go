package jenkins

import (
	"context"

	"git.4thena.io/4thena/abys/internal/model"
)

type JenkinsCi struct{}

func NewJenkinsCi(url, token string) (*JenkinsCi, error) {
	return &JenkinsCi{}, nil
}

func (c *JenkinsCi) ActivateRepo(ctx context.Context, ID int64) (*model.CIRepo, error) {
	return &model.CIRepo{
		ID: 1,
		URL: "https://jenkins.com/1",
	}, nil
}

func (c *JenkinsCi) GetBuilds(ctx context.Context, ID int64) ([]model.Build, error) {
	res := make([]model.Build, 1)
	res[0] = model.Build{}
	return res, nil
}
