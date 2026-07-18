package service

import (
	"context"
	"errors"
	"testing"

	"github.com/4thena-io/abyss/internal/integration/forge/mocks"
	"github.com/4thena-io/abyss/internal/model"
	svcmocks "github.com/4thena-io/abyss/internal/service/mocks"
	"github.com/stretchr/testify/mock"
)

func TestTemplateService_SaveTemplate(t *testing.T) {
	t.Run("returns ErrConflict when name already exists", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByName(context.Background(), "taken").Return(&model.Template{ID: 1, Name: "taken"}, nil)

		svc := NewTemplateService(repo, nil, "owner")

		_, err := svc.SaveTemplate(context.Background(), &model.Template{Name: "taken"})
		if !errors.Is(err, ErrConflict) {
			t.Fatalf("expected ErrConflict, got %v", err)
		}
	})

	t.Run("saves when name is free", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		template := &model.Template{Name: "new"}
		repo.EXPECT().GetByName(context.Background(), "new").Return(nil, nil)
		repo.EXPECT().Save(context.Background(), template).Return(nil)

		svc := NewTemplateService(repo, nil, "owner")

		got, err := svc.SaveTemplate(context.Background(), template)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != template {
			t.Fatalf("expected the same template instance to be returned")
		}
	})
}

func TestTemplateService_UpdateTemplate(t *testing.T) {
	t.Run("returns ErrNotFound when template missing", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByID(context.Background(), uint(1)).Return(nil, nil)

		svc := NewTemplateService(repo, nil, "owner")

		_, err := svc.UpdateTemplate(context.Background(), 1, 1, false, "n", "d", "k", "l")
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden for non-admin non-creator", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByID(context.Background(), uint(1)).Return(&model.Template{ID: 1, CreatorID: 42}, nil)

		svc := NewTemplateService(repo, nil, "owner")

		_, err := svc.UpdateTemplate(context.Background(), 1, 999, false, "n", "d", "k", "l")
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("updates fields for creator", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		template := &model.Template{ID: 1, CreatorID: 42, Name: "old"}
		repo.EXPECT().GetByID(context.Background(), uint(1)).Return(template, nil)
		repo.EXPECT().Update(context.Background(), template).Return(nil)

		svc := NewTemplateService(repo, nil, "owner")

		got, err := svc.UpdateTemplate(context.Background(), 1, 42, false, "renamed", "new desc", "service", "go")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "renamed" || got.Description != "new desc" || got.Kind != "service" || got.Language != "go" {
			t.Fatalf("template not updated as expected: %+v", got)
		}
	})
}

func TestTemplateService_CreateBlankTemplate(t *testing.T) {
	t.Run("errors when forge is not configured", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		svc := NewTemplateService(repo, nil, "owner")

		_, err := svc.CreateBlankTemplate(context.Background(), 1, "name", "desc", "service", "go")
		if err == nil {
			t.Fatal("expected an error when forge is not configured")
		}
	})

	t.Run("returns ErrConflict when name already exists", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByName(context.Background(), "taken").Return(&model.Template{ID: 1}, nil)
		forge := mocks.NewMockForge(t)

		svc := NewTemplateService(repo, forge, "owner")

		_, err := svc.CreateBlankTemplate(context.Background(), 1, "taken", "desc", "service", "go")
		if !errors.Is(err, ErrConflict) {
			t.Fatalf("expected ErrConflict, got %v", err)
		}
	})

	t.Run("rolls back the forge repo when save fails", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByName(context.Background(), "name").Return(nil, nil)
		repo.On("Save", mock.Anything, mock.AnythingOfType("*model.Template")).Return(errBoom)

		forge := mocks.NewMockForge(t)
		forge.EXPECT().CreateRepo(context.Background(), "owner", "name").Return(&model.Repo{ID: 1, URL: "u", CloneURL: "c"}, nil)
		forge.EXPECT().DeleteRepo(context.Background(), "owner", "name").Return(nil)

		svc := NewTemplateService(repo, forge, "owner")

		_, err := svc.CreateBlankTemplate(context.Background(), 1, "name", "desc", "service", "go")
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})

	t.Run("creates repo and saves template", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByName(context.Background(), "name").Return(nil, nil)
		repo.On("Save", mock.Anything, mock.AnythingOfType("*model.Template")).Return(nil)

		forge := mocks.NewMockForge(t)
		forge.EXPECT().CreateRepo(context.Background(), "owner", "name").Return(&model.Repo{ID: 1, URL: "https://forge/owner/name", CloneURL: "https://forge/owner/name.git"}, nil)

		svc := NewTemplateService(repo, forge, "owner")

		got, err := svc.CreateBlankTemplate(context.Background(), 7, "name", "desc", "service", "go")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.RepoURL != "https://forge/owner/name" || got.CreatorID != 7 {
			t.Fatalf("unexpected template: %+v", got)
		}
	})
}

func TestTemplateService_DeleteTemplate(t *testing.T) {
	t.Run("returns ErrNotFound when template missing", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByID(context.Background(), uint(1)).Return(nil, nil)

		svc := NewTemplateService(repo, nil, "owner")

		err := svc.DeleteTemplate(context.Background(), 1, 1, false)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden for non-admin non-creator", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		repo.EXPECT().GetByID(context.Background(), uint(1)).Return(&model.Template{ID: 1, CreatorID: 42}, nil)

		svc := NewTemplateService(repo, nil, "owner")

		err := svc.DeleteTemplate(context.Background(), 1, 999, false)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("deletes when caller is admin", func(t *testing.T) {
		repo := svcmocks.NewMockTemplateRepository(t)
		template := &model.Template{ID: 1, CreatorID: 42}
		repo.EXPECT().GetByID(context.Background(), uint(1)).Return(template, nil)
		repo.EXPECT().Delete(context.Background(), template).Return(nil)

		svc := NewTemplateService(repo, nil, "owner")

		if err := svc.DeleteTemplate(context.Background(), 1, 999, true); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
