package jenkins

import (
	"context"

	"git.d4ramirez.com/project-abyss/abys-api/internal/dto"
)

type JenkinsCi struct{}

func NewJenkinsCi(url, token string) (*JenkinsCi, error) {
	return &JenkinsCi{}, nil
}

func (c *JenkinsCi) ActivateRepo(ctx context.Context, req dto.ActivateRepoRequestDTO) (*dto.ActivateRepoResponseDTO, error) {
	return &dto.ActivateRepoResponseDTO{
		RepoId: 1,
		RepoUrl: "https://jenkins.com/1",
	}, nil
}
