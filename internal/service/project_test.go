package service

import (
	"context"
	"errors"
	"testing"

	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service/mocks"
	"github.com/stretchr/testify/mock"
)

func TestProjectService_SaveProject(t *testing.T) {
	t.Run("returns ErrConflict when name already taken", func(t *testing.T) {
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetByName(mock.Anything, "taken").Return(&model.Project{ID: 1, Name: "taken"}, nil)
		svc := NewProjectService(repo)

		_, err := svc.SaveProject(context.Background(), &model.Project{Name: "taken"})
		if !errors.Is(err, ErrConflict) {
			t.Fatalf("expected ErrConflict, got %v", err)
		}
	})

	t.Run("propagates lookup error", func(t *testing.T) {
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetByName(mock.Anything, "x").Return(nil, errBoom)
		svc := NewProjectService(repo)

		_, err := svc.SaveProject(context.Background(), &model.Project{Name: "x"})
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})

	t.Run("saves when name is free", func(t *testing.T) {
		project := &model.Project{Name: "new-project"}
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetByName(mock.Anything, "new-project").Return(nil, nil)
		repo.EXPECT().Save(mock.Anything, project).Return(nil)
		svc := NewProjectService(repo)

		got, err := svc.SaveProject(context.Background(), project)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != project {
			t.Fatalf("expected the same project instance to be returned")
		}
	})
}

func TestProjectService_UpdateProject(t *testing.T) {
	t.Run("returns ErrNotFound when project missing", func(t *testing.T) {
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		svc := NewProjectService(repo)

		_, err := svc.UpdateProject(context.Background(), 1, 42, false, "new", "desc", nil)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden for non-admin non-creator", func(t *testing.T) {
		existing := &model.Project{ID: 1, Name: "old", CreatorID: 42}
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(existing, nil)
		svc := NewProjectService(repo)

		_, err := svc.UpdateProject(context.Background(), 1, 999, false, "new", "desc", nil)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("allows admin to update someone else's project", func(t *testing.T) {
		p := &model.Project{ID: 1, Name: "old", CreatorID: 42}
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(p, nil)
		repo.EXPECT().Update(mock.Anything, p).Return(nil)
		svc := NewProjectService(repo)

		teamID := uint(7)
		got, err := svc.UpdateProject(context.Background(), 1, 999, true, "renamed", "new desc", &teamID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "renamed" || got.Description != "new desc" || got.TeamID == nil || *got.TeamID != 7 {
			t.Fatalf("project was not updated as expected: %+v", got)
		}
	})

	t.Run("allows creator to update own project", func(t *testing.T) {
		p := &model.Project{ID: 1, Name: "old", CreatorID: 42}
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(p, nil)
		repo.EXPECT().Update(mock.Anything, p).Return(nil)
		svc := NewProjectService(repo)

		got, err := svc.UpdateProject(context.Background(), 1, 42, false, "renamed", "new desc", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "renamed" {
			t.Fatalf("expected name to be updated, got %q", got.Name)
		}
	})

	t.Run("propagates update error", func(t *testing.T) {
		p := &model.Project{ID: 1, Name: "old", CreatorID: 42}
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(p, nil)
		repo.EXPECT().Update(mock.Anything, p).Return(errBoom)
		svc := NewProjectService(repo)

		_, err := svc.UpdateProject(context.Background(), 1, 42, false, "renamed", "desc", nil)
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})
}

func TestProjectService_DeleteProject(t *testing.T) {
	t.Run("returns ErrNotFound when project missing", func(t *testing.T) {
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		svc := NewProjectService(repo)

		err := svc.DeleteProject(context.Background(), 1, 42, false)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden for non-admin non-creator", func(t *testing.T) {
		p := &model.Project{ID: 1, CreatorID: 42}
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(p, nil)
		svc := NewProjectService(repo)

		err := svc.DeleteProject(context.Background(), 1, 999, false)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("deletes when caller is creator", func(t *testing.T) {
		p := &model.Project{ID: 1, CreatorID: 42}
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(p, nil)
		repo.EXPECT().Delete(mock.Anything, p).Return(nil)
		svc := NewProjectService(repo)

		if err := svc.DeleteProject(context.Background(), 1, 42, false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestProjectService_GetAllProjects(t *testing.T) {
	t.Run("returns projects from repository", func(t *testing.T) {
		want := []model.Project{{ID: 1}, {ID: 2}}
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetAll(mock.Anything).Return(want, nil)
		svc := NewProjectService(repo)

		got, err := svc.GetAllProjects(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("expected 2 projects, got %d", len(got))
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repo := mocks.NewMockProjectRepository(t)
		repo.EXPECT().GetAll(mock.Anything).Return(nil, errBoom)
		svc := NewProjectService(repo)

		_, err := svc.GetAllProjects(context.Background())
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})
}

func TestProjectService_GetProjectByID(t *testing.T) {
	repo := mocks.NewMockProjectRepository(t)
	repo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.Project{ID: 1}, nil)
	repo.EXPECT().GetByID(mock.Anything, uint(2)).Return(nil, nil)
	svc := NewProjectService(repo)

	got, err := svc.GetProjectByID(context.Background(), 1)
	if err != nil || got == nil {
		t.Fatalf("expected project, got %+v, err %v", got, err)
	}

	got, err = svc.GetProjectByID(context.Background(), 2)
	if err != nil || got != nil {
		t.Fatalf("expected nil project, got %+v, err %v", got, err)
	}
}
