package repository

import (
	"context"

	"git.4thena.io/4thena/abys/internal/model"
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

func (r *AppRepository) SaveApp(ctx context.Context, app *model.App) error {
	result := r.db.WithContext(ctx).Create(app)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AppRepository) GetAllApps(ctx context.Context) ([]model.App, error) {
	var applications []model.App
	result := r.db.WithContext(ctx).Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

func (r *AppRepository) GetAppById(ctx context.Context, id uint) (*model.App, error) {
	var application model.App
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&application)
	if result.Error != nil {
		return nil, result.Error
	}
	return &application, nil
}

func (r *AppRepository) GetAppByName(ctx context.Context, name string) (*model.App, error) {
	var application model.App
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&application)
	if result.Error != nil {
		return nil, result.Error
	}
	return &application, nil
}

func (r *AppRepository) GetAppsByProject(ctx context.Context, id uint) ([]model.App, error) {
	var applications []model.App
	result := r.db.WithContext(ctx).Where("project_id = ?", id).Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

func (r *AppRepository) DeleteApp(ctx context.Context, app *model.App) error {
	return r.db.WithContext(ctx).Delete(&app).Error
}
