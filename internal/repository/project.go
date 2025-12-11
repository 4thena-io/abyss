package repository

import (
	"context"
	"time"

	"git.4thena.io/4thena/abys/internal/constant"
	"git.4thena.io/4thena/abys/internal/dto"
	"git.4thena.io/4thena/abys/internal/model"
	"github.com/google/uuid"
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

func (r *ProjectRepository) SaveProject(ctx context.Context, project *dto.CreateProjectRecordDto) (*model.Project, error) {
	newProject := model.Project{
		Id:          uuid.NewString(),
		Name:        project.Name,
		Description: project.Description,
		Status:      constant.StatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := r.db.WithContext(ctx).Create(newProject)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newProject, nil
}

func (r *ProjectRepository) GetAllProjects(ctx context.Context) ([]model.Project, error) {
	var projects []model.Project
	result := r.db.WithContext(ctx).Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

func (r *ProjectRepository) GetProjectById(ctx context.Context, id string) (*model.Project, error) {
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
