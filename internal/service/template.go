package service

import (
	"context"
	"fmt"

	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/repository"
	"gorm.io/gorm"
)

type TemplateService struct {
	repository repository.TemplateRepository
}

func NewTemplateService(repository repository.TemplateRepository) *TemplateService {
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
	data, err := s.repository.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("template not found")
	}

	return data, nil
}

func (s *TemplateService) GetTemplateByID(ctx context.Context, id uint) (*model.Template, error) {
	template, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, fmt.Errorf("template not found")
	}

	return template, nil
}

func (s *TemplateService) SaveTemplate(ctx context.Context, template *model.Template) (*model.Template, error) {
	existing, err := s.repository.GetByName(ctx, template.Name)
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("the template already exists")
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
		return fmt.Errorf("template doesn't exist")
	}
	return s.repository.Delete(ctx, template)
}
