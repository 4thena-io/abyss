package service

import (
	"context"
	"fmt"

	"git.4thena.io/4thena/abys/internal/model"
	"git.4thena.io/4thena/abys/internal/repository"
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
	data, err := s.repository.GetAllTemplates(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *TemplateService) GetTemplateByName(ctx context.Context, name string) (*model.Template, error) {
	data, err := s.repository.GetTemplateByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("template not found")
	}

	return data, nil
}

func (s *TemplateService) SaveTemplate(ctx context.Context, template *model.Template) (*model.Template, error) {
	existing, err := s.repository.GetTemplateByName(ctx, template.Name)
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("the template already exists")
	}

	err = s.repository.SaveTemplate(ctx, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}
