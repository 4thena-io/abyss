package service

import (
	"context"
	"errors"
	"testing"

	"github.com/4thena-io/abyss/internal/model"
)

func TestProjectService_SaveProject(t *testing.T) {
	t.Run("returns ErrConflict when name already taken", func(t *testing.T) {
		repo := &mockProjectRepository{
			GetByNameFn: func(ctx context.Context, name string) (*model.Project, error) {
				return &model.Project{ID: 1, Name: name}, nil
			},
		}
		svc := NewProjectService(repo)

		_, err := svc.SaveProject(context.Background(), &model.Project{Name: "taken"})
		if !errors.Is(err, ErrConflict) {
			t.Fatalf("expected ErrConflict, got %v", err)
		}
	})

	t.Run("propagates lookup error", func(t *testing.T) {
		repo := &mockProjectRepository{
			GetByNameFn: func(ctx context.Context, name string) (*model.Project, error) {
				return nil, errBoom
			},
		}
		svc := NewProjectService(repo)

		_, err := svc.SaveProject(context.Background(), &model.Project{Name: "x"})
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})

	t.Run("saves when name is free", func(t *testing.T) {
		var saved *model.Project
		repo := &mockProjectRepository{
			GetByNameFn: func(ctx context.Context, name string) (*model.Project, error) { return nil, nil },
			SaveFn: func(ctx context.Context, project *model.Project) error {
				saved = project
				return nil
			},
		}
		svc := NewProjectService(repo)

		project := &model.Project{Name: "new-project"}
		got, err := svc.SaveProject(context.Background(), project)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != project || saved != project {
			t.Fatalf("expected the same project instance to be returned and saved")
		}
	})
}

func TestProjectService_UpdateProject(t *testing.T) {
	existing := &model.Project{ID: 1, Name: "old", CreatorID: 42}

	t.Run("returns ErrNotFound when project missing", func(t *testing.T) {
		repo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return nil, nil },
		}
		svc := NewProjectService(repo)

		_, err := svc.UpdateProject(context.Background(), 1, 42, false, "new", "desc", nil)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden for non-admin non-creator", func(t *testing.T) {
		repo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return existing, nil },
		}
		svc := NewProjectService(repo)

		_, err := svc.UpdateProject(context.Background(), 1, 999, false, "new", "desc", nil)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("allows admin to update someone else's project", func(t *testing.T) {
		p := &model.Project{ID: 1, Name: "old", CreatorID: 42}
		repo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return p, nil },
			UpdateFn:  func(ctx context.Context, project *model.Project) error { return nil },
		}
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
		repo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return p, nil },
			UpdateFn:  func(ctx context.Context, project *model.Project) error { return nil },
		}
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
		repo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return p, nil },
			UpdateFn:  func(ctx context.Context, project *model.Project) error { return errBoom },
		}
		svc := NewProjectService(repo)

		_, err := svc.UpdateProject(context.Background(), 1, 42, false, "renamed", "desc", nil)
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})
}

func TestProjectService_DeleteProject(t *testing.T) {
	t.Run("returns ErrNotFound when project missing", func(t *testing.T) {
		repo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return nil, nil },
		}
		svc := NewProjectService(repo)

		err := svc.DeleteProject(context.Background(), 1, 42, false)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden for non-admin non-creator", func(t *testing.T) {
		p := &model.Project{ID: 1, CreatorID: 42}
		repo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return p, nil },
		}
		svc := NewProjectService(repo)

		err := svc.DeleteProject(context.Background(), 1, 999, false)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("deletes when caller is creator", func(t *testing.T) {
		p := &model.Project{ID: 1, CreatorID: 42}
		deleted := false
		repo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return p, nil },
			DeleteFn: func(ctx context.Context, project *model.Project) error {
				deleted = true
				return nil
			},
		}
		svc := NewProjectService(repo)

		if err := svc.DeleteProject(context.Background(), 1, 42, false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !deleted {
			t.Fatalf("expected repository Delete to be called")
		}
	})
}

func TestProjectService_GetAllProjects(t *testing.T) {
	t.Run("returns projects from repository", func(t *testing.T) {
		want := []model.Project{{ID: 1}, {ID: 2}}
		repo := &mockProjectRepository{
			GetAllFn: func(ctx context.Context) ([]model.Project, error) { return want, nil },
		}
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
		repo := &mockProjectRepository{
			GetAllFn: func(ctx context.Context) ([]model.Project, error) { return nil, errBoom },
		}
		svc := NewProjectService(repo)

		_, err := svc.GetAllProjects(context.Background())
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})
}

func TestProjectService_GetProjectByID(t *testing.T) {
	repo := &mockProjectRepository{
		GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) {
			if id == 1 {
				return &model.Project{ID: 1}, nil
			}
			return nil, nil
		},
	}
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
