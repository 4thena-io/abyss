package repository

import (
	"context"
	"errors"

	"github.com/4thena-io/abyss/internal/model"
	"gorm.io/gorm"
)

type AppRepository struct {
	db *gorm.DB
}

func NewAppRepository(db *gorm.DB) *AppRepository {
	return &AppRepository{
		db: db,
	}
}

func (r *AppRepository) Save(ctx context.Context, app *model.App) error {
	result := r.db.WithContext(ctx).Create(app)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AppRepository) GetAll(ctx context.Context) ([]model.App, error) {
	var applications []model.App
	result := r.db.WithContext(ctx).Preload("Creator").Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

func (r *AppRepository) GetByID(ctx context.Context, id uint) (*model.App, error) {
	var application model.App
	result := r.db.WithContext(ctx).Preload("Creator").Where("id = ?", id).First(&application)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &application, nil
}

func (r *AppRepository) GetByName(ctx context.Context, name string) (*model.App, error) {
	var application model.App
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&application)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &application, nil
}

func (r *AppRepository) GetByProject(ctx context.Context, id uint) ([]model.App, error) {
	var applications []model.App
	result := r.db.WithContext(ctx).Where("project_id = ?", id).Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

func (r *AppRepository) GetByTemplate(ctx context.Context, id uint) ([]model.App, error) {
	var applications []model.App
	result := r.db.WithContext(ctx).Where("template_id = ?", id).Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

func (r *AppRepository) Update(ctx context.Context, app *model.App) error {
	return r.db.WithContext(ctx).Save(app).Error
}

func (r *AppRepository) Delete(ctx context.Context, app *model.App) error {
	return r.db.WithContext(ctx).Delete(&app).Error
}

func (r *AppRepository) CountByTeam(ctx context.Context, teamID uint) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.App{}).
		Joins("JOIN projects ON projects.id = apps.project_id").
		Where("projects.team_id = ?", teamID).
		Count(&count)
	return count, result.Error
}
