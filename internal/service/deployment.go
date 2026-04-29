package service

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
)

type DeploymentRepository interface {
	Save(ctx context.Context, deployment *model.Deployment) error
	GetByApp(ctx context.Context, appID uint) ([]model.Deployment, error)
}

type DeploymentService struct {
	deploymentRepository DeploymentRepository
}

func NewDeploymentService(deploymentRepository DeploymentRepository) *DeploymentService {
	return &DeploymentService{deploymentRepository: deploymentRepository}
}


func (s *DeploymentService) GetDeploymentsByApp(ctx context.Context, appID uint) ([]model.Deployment, error) {
	return s.deploymentRepository.GetByApp(ctx, appID)
}

func (s *DeploymentService) SaveDeployment(ctx context.Context, deployment *model.Deployment) error {
	return s.deploymentRepository.Save(ctx, deployment)
}
