package repository

import (
	"context"

	"git.4thena.io/4thena/abys/internal/model"
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

func (r *TemplateRepository) SaveTemplate(ctx context.Context, template *model.Template) error {
	result := r.db.WithContext(ctx).Create(template)
	if result.Error != nil {
		return result.Error
	}
	return nil
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
