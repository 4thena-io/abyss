package service

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/repository"
)

type DeploymentService struct {
	deploymentRepository repository.DeploymentRepository
}

func NewDeploymentService(deploymentRepository repository.DeploymentRepository) *DeploymentService {
	return &DeploymentService{deploymentRepository: deploymentRepository}
}


func (s *DeploymentService) GetDeploymentsByApp(ctx context.Context, appID uint) ([]model.Deployment, error) {
	return s.deploymentRepository.GetByApp(ctx, appID)
}

func (s *DeploymentService) SaveDeployment(ctx context.Context, deployment *model.Deployment) error {
	return s.deploymentRepository.Save(ctx, deployment)
}
