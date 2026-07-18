package service

import (
	"context"
	"fmt"

	"github.com/4thena-io/abyss/internal/integration/forge"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/rs/zerolog/log"
)

type TemplateRepository interface {
	Save(ctx context.Context, template *model.Template) error
	Update(ctx context.Context, template *model.Template) error
	GetAll(ctx context.Context) ([]model.Template, error)
	GetByID(ctx context.Context, id uint) (*model.Template, error)
	GetByName(ctx context.Context, name string) (*model.Template, error)
	Delete(ctx context.Context, template *model.Template) error
}

type TemplateService struct {
	repository TemplateRepository
	forge      forge.Forge
	owner      string
}

func NewTemplateService(repository TemplateRepository, forgeProvider forge.Forge, owner string) *TemplateService {
	return &TemplateService{
		repository: repository,
		forge:      forgeProvider,
		owner:      owner,
	}
}

func (s *TemplateService) GetAllTemplates(ctx context.Context) ([]model.Template, error) {
	return s.repository.GetAll(ctx)
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
	if err := s.repository.Save(ctx, template); err != nil {
		return nil, err
	}
	return template, nil
}

func (s *TemplateService) UpdateTemplate(ctx context.Context, id, callerID uint, isAdmin bool, name, description, kind, language string) (*model.Template, error) {
	t, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrNotFound
	}
	if !isAdmin && t.CreatorID != callerID {
		return nil, ErrForbidden
	}
	t.Name = name
	t.Description = description
	t.Kind = kind
	t.Language = language
	if err := s.repository.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TemplateService) CreateBlankTemplate(ctx context.Context, creatorID uint, name, description, kind, language string) (*model.Template, error) {
	if s.forge == nil {
		return nil, fmt.Errorf("forge not configured")
	}

	existing, err := s.repository.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrConflict
	}

	repo, err := s.forge.CreateRepo(ctx, s.owner, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create repo: %w", err)
	}

	template := &model.Template{
		Name:        name,
		Description: description,
		Kind:        kind,
		Language:    language,
		RepoURL:     repo.URL,
		CloneURL:    repo.CloneURL,
		CreatorID:   creatorID,
	}
	if err := s.repository.Save(ctx, template); err != nil {
		if delErr := s.forge.DeleteRepo(ctx, s.owner, name); delErr != nil {
			log.Warn().Err(delErr).Str("repo", name).Msg("failed to roll back forge repo after template save failure")
		}
		return nil, err
	}
	return template, nil
}

func (s *TemplateService) DeleteTemplate(ctx context.Context, id, callerID uint, isAdmin bool) error {
	template, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if template == nil {
		return ErrNotFound
	}
	if !isAdmin && template.CreatorID != callerID {
		return ErrForbidden
	}
	return s.repository.Delete(ctx, template)
}
