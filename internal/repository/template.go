package repository

import (
	"context"
	"time"

	"git.4thena.io/4thena/abys/internal/constant"
	"git.4thena.io/4thena/abys/internal/dto/request"
	"git.4thena.io/4thena/abys/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{
		db: db,
	}
}

func (r *TemplateRepository) SaveTemplate(ctx context.Context, app *request.CreateTemplate) (*model.Template, error) {
	newTemplate := model.Template{
		Id:          uuid.NewString(),
		Name:        app.Name,
		Description: app.Description,
		Kind:        app.Kind,
		Language:    app.Language,
		RepoUrl:     app.RepoUrl,
		CloneUrl:    app.CloneUrl,
		Status:      constant.StatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := r.db.WithContext(ctx).Create(newTemplate)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newTemplate, nil
}

func (r *TemplateRepository) GetAllTemplates(ctx context.Context) ([]model.Template, error) {
	var templates []model.Template
	result := r.db.WithContext(ctx).Find(&templates)
	if result.Error != nil {
		return nil, result.Error
	}
	return templates, nil
}

func (r *TemplateRepository) GetTemplateById(ctx context.Context, id string) (*model.Template, error) {
	var template model.Template
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&template)
	if result.Error != nil {
		return nil, result.Error
	}
	return &template, nil
}

func (r *TemplateRepository) GetTemplateByName(ctx context.Context, name string) (*model.Template, error) {
	var template model.Template
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&template)
	if result.Error != nil {
		return nil, result.Error
	}
	return &template, nil
}

func (r *TemplateRepository) DeleteTemplate(ctx context.Context, template *model.Template) error {
	return r.db.WithContext(ctx).Delete(&template).Error
}
