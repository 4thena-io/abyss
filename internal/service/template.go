package service

import (
	"context"

	"github.com/4thena-io/abyss/internal/model"
)

type TemplateRepository interface {
	Save(ctx context.Context, template *model.Template) error
	GetAll(ctx context.Context) ([]model.Template, error)
	GetByID(ctx context.Context, id uint) (*model.Template, error)
	GetByName(ctx context.Context, name string) (*model.Template, error)
	Delete(ctx context.Context, template *model.Template) error
}

type TemplateService struct {
	repository TemplateRepository
}

func NewTemplateService(repository TemplateRepository) *TemplateService {
	return &TemplateService{
		repository: repository,
	}
}

func (s *TemplateService) GetAllTemplates(ctx context.Context) ([]model.Template, error) {
	data, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *TemplateService) GetTemplateByName(ctx context.Context, name string) (*model.Template, error) {
	return s.repository.GetByName(ctx, name)
}

func (s *TemplateService) GetTemplateByID(ctx context.Context, id uint) (*model.Template, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *TemplateService) SaveTemplate(ctx context.Context, template *model.Template) (*model.Template, error) {
	existing, err := s.repository.GetByName(ctx, template.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrConflict
	}

	err = s.repository.Save(ctx, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (s *TemplateService) DeleteTemplate(ctx context.Context, id uint) error {
	template, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if template == nil {
		return ErrNotFound
	}
	return s.repository.Delete(ctx, template)
}
