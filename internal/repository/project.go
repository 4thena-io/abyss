package repository

import (
	"context"
	"errors"

	"github.com/4thena-io/abyss/internal/model"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

func (r *ProjectRepository) Save(ctx context.Context, project *model.Project) error {
	result := r.db.WithContext(ctx).Create(project)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProjectRepository) GetAll(ctx context.Context) ([]model.Project, error) {
	var projects []model.Project
	result := r.db.WithContext(ctx).Preload("Creator").Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uint) (*model.Project, error) {
	var project model.Project
	result := r.db.WithContext(ctx).Preload("Creator").Where("id = ?", id).First(&project)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (r *ProjectRepository) GetByName(ctx context.Context, name string) (*model.Project, error) {
	var project model.Project
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&project)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *ProjectRepository) Delete(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Delete(&project).Error
}

func (r *ProjectRepository) GetByTeam(ctx context.Context, teamID uint) ([]model.Project, error) {
	var projects []model.Project
	result := r.db.WithContext(ctx).Where("team_id = ?", teamID).Find(&projects)
	return projects, result.Error
}

func (r *ProjectRepository) CountByTeam(ctx context.Context, teamID uint) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Project{}).Where("team_id = ?", teamID).Count(&count)
	return count, result.Error
}
