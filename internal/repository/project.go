package repository

import (
	"context"

	"git.4thena.io/4thena/abys/internal/model"
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

func (r *ProjectRepository) SaveProject(ctx context.Context, project *model.Project) error {
	result := r.db.WithContext(ctx).Create(project)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ProjectRepository) GetAllProjects(ctx context.Context) ([]model.Project, error) {
	var projects []model.Project
	result := r.db.WithContext(ctx).Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

func (r *ProjectRepository) GetProjectByID(ctx context.Context, id uint) (*model.Project, error) {
	var project model.Project
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&project)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (r *ProjectRepository) GetProjectByName(ctx context.Context, name string) (*model.Project, error) {
	var project model.Project
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&project)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}
