package repository

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
	"gorm.io/gorm"
)

type DeploymentRepository struct {
	db *gorm.DB
}

func NewDeploymentRepository(db *gorm.DB) *DeploymentRepository {
	return &DeploymentRepository{db: db}
}

func (r *DeploymentRepository) Save(ctx context.Context, deployment *model.Deployment) error {
	return r.db.WithContext(ctx).Create(deployment).Error
}

func (r *DeploymentRepository) GetByApp(ctx context.Context, appID uint) ([]model.Deployment, error) {
	var deployments []model.Deployment
	result := r.db.WithContext(ctx).Where("app_id = ?", appID).Order("deployed_at DESC").Find(&deployments)
	return deployments, result.Error
}

func (r *DeploymentRepository) GetByID(ctx context.Context, id uint) (*model.Deployment, error) {
	var deployment model.Deployment
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&deployment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &deployment, nil
}

func (r *DeploymentRepository) Delete(ctx context.Context, deployment *model.Deployment) error {
	return r.db.WithContext(ctx).Delete(deployment).Error
}
