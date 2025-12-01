package repository

import (
	"context"
	"time"

	"git.d4ramirez.com/project-abyss/abys-api/internal/constant"
	"git.d4ramirez.com/project-abyss/abys-api/internal/dto"
	"git.d4ramirez.com/project-abyss/abys-api/internal/model"
	"github.com/google/uuid"
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

func (r *AppRepository) SaveApp(ctx context.Context, app *dto.CreateAppRecordDto) (*model.App, error) {
	newApp := model.App{
		Id:          uuid.NewString(),
		Name:        app.Name,
		Description: app.Description,
		RepoId:      app.RepoId,
		RepoUrl:     app.RepoUrl,
		CloneUrl:    app.CloneUrl,
		CiId:        app.CiId,
		CiUrl:       app.CiUrl,
		Project:     app.Project,
		Status:      constant.StatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := r.db.WithContext(ctx).Create(newApp)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newApp, nil
}

func (r *AppRepository) GetAllApps(ctx context.Context) ([]model.App, error) {
	var applications []model.App
	result := r.db.WithContext(ctx).Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

func (r *AppRepository) GetAppById(ctx context.Context, id string) (*model.App, error) {
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

func (r *AppRepository) DeleteApp(ctx context.Context, app *model.App) error {
	return r.db.WithContext(ctx).Delete(&app).Error
}
