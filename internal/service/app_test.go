package service

import (
	"context"
	"errors"
	"testing"

	"github.com/4thena-io/abyss/internal/model"
)

func newTestAppService(appRepo AppRepository, projectRepo ProjectRepository, templateRepo TemplateRepository, forge *mockForge, ci *mockCI) *AppService {
	return NewAppService(appRepo, projectRepo, templateRepo, forge, ci, nil, "owner", "https://abyss.example", "secret", "main")
}

func TestAppService_CreateAppFromTemplate(t *testing.T) {
	t.Run("errors when project id set but project missing", func(t *testing.T) {
		projectID := uint(5)
		projectRepo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return nil, nil },
		}
		svc := newTestAppService(&mockAppRepository{}, projectRepo, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		_, err := svc.CreateAppFromTemplate(context.Background(), &model.App{Name: "app", ProjectID: &projectID, TemplateID: ptr(uint(1))})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("errors when template id is nil", func(t *testing.T) {
		svc := newTestAppService(&mockAppRepository{}, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		_, err := svc.CreateAppFromTemplate(context.Background(), &model.App{Name: "app"})
		if err == nil {
			t.Fatal("expected error for missing template ID")
		}
	})

	t.Run("errors when template does not exist", func(t *testing.T) {
		templateRepo := &mockTemplateRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Template, error) { return nil, nil },
		}
		svc := newTestAppService(&mockAppRepository{}, &mockProjectRepository{}, templateRepo, &mockForge{}, &mockCI{})

		_, err := svc.CreateAppFromTemplate(context.Background(), &model.App{Name: "app", TemplateID: ptr(uint(1))})
		if err == nil {
			t.Fatal("expected error for missing template")
		}
	})

	t.Run("returns ErrConflict when app name already exists", func(t *testing.T) {
		templateRepo := &mockTemplateRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Template, error) {
				return &model.Template{ID: 1, CloneURL: "https://forge/template.git"}, nil
			},
		}
		appRepo := &mockAppRepository{
			GetByNameFn: func(ctx context.Context, name string) (*model.App, error) {
				return &model.App{ID: 1, Name: name}, nil
			},
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, templateRepo, &mockForge{}, &mockCI{})

		_, err := svc.CreateAppFromTemplate(context.Background(), &model.App{Name: "taken", TemplateID: ptr(uint(1))})
		if !errors.Is(err, ErrConflict) {
			t.Fatalf("expected ErrConflict, got %v", err)
		}
	})

	t.Run("propagates forge repo creation failure without touching git", func(t *testing.T) {
		templateRepo := &mockTemplateRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Template, error) {
				return &model.Template{ID: 1, CloneURL: "https://forge/template.git"}, nil
			},
		}
		appRepo := &mockAppRepository{
			GetByNameFn: func(ctx context.Context, name string) (*model.App, error) { return nil, nil },
		}
		forge := &mockForge{
			CreateRepoFn: func(ctx context.Context, owner, name string) (*model.Repo, error) {
				return nil, errBoom
			},
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, templateRepo, forge, &mockCI{})

		_, err := svc.CreateAppFromTemplate(context.Background(), &model.App{Name: "app", TemplateID: ptr(uint(1))})
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})
}

func TestAppService_CreateAppFromRepo(t *testing.T) {
	t.Run("errors when project id set but project missing", func(t *testing.T) {
		projectID := uint(5)
		projectRepo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return nil, nil },
		}
		svc := newTestAppService(&mockAppRepository{}, projectRepo, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		_, err := svc.CreateAppFromRepo(context.Background(), &model.App{Name: "app", ProjectID: &projectID})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns ErrConflict when app name already exists", func(t *testing.T) {
		appRepo := &mockAppRepository{
			GetByNameFn: func(ctx context.Context, name string) (*model.App, error) {
				return &model.App{ID: 1, Name: name}, nil
			},
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		_, err := svc.CreateAppFromRepo(context.Background(), &model.App{Name: "taken", RepoID: 99})
		if !errors.Is(err, ErrConflict) {
			t.Fatalf("expected ErrConflict, got %v", err)
		}
	})

	t.Run("errors when repo does not exist in forge", func(t *testing.T) {
		appRepo := &mockAppRepository{
			GetByNameFn: func(ctx context.Context, name string) (*model.App, error) { return nil, nil },
		}
		forge := &mockForge{
			GetRepoFn: func(ctx context.Context, id int64) (*model.Repo, error) { return nil, errBoom },
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, forge, &mockCI{})

		_, err := svc.CreateAppFromRepo(context.Background(), &model.App{Name: "app", RepoID: 99})
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})
}

func TestAppService_GetAppByID(t *testing.T) {
	appRepo := &mockAppRepository{
		GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) {
			if id == 1 {
				return &model.App{ID: 1}, nil
			}
			return nil, nil
		},
	}
	svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

	got, err := svc.GetAppByID(context.Background(), 1)
	if err != nil || got == nil {
		t.Fatalf("expected app, got %+v, err %v", got, err)
	}
}

func TestAppService_GetAppBuilds(t *testing.T) {
	t.Run("returns ErrNotFound when app missing", func(t *testing.T) {
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return nil, nil },
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		_, err := svc.GetAppBuilds(context.Background(), 1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns builds from ci", func(t *testing.T) {
		app := &model.App{ID: 1, CIID: 55, CISlug: "owner/app"}
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return app, nil },
		}
		ci := &mockCI{
			GetBuildsFn: func(ctx context.Context, ciRepoID int64, slug string) ([]model.Build, error) {
				if ciRepoID != 55 || slug != "owner/app" {
					t.Fatalf("unexpected args: %d %s", ciRepoID, slug)
				}
				return []model.Build{{ID: 1}}, nil
			},
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, ci)

		builds, err := svc.GetAppBuilds(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(builds) != 1 {
			t.Fatalf("expected 1 build, got %d", len(builds))
		}
	})
}

func TestAppService_GetAppsByProject(t *testing.T) {
	t.Run("returns ErrNotFound when project missing", func(t *testing.T) {
		projectRepo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return nil, nil },
		}
		svc := newTestAppService(&mockAppRepository{}, projectRepo, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		_, err := svc.GetAppsByProject(context.Background(), 1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns apps for project", func(t *testing.T) {
		projectRepo := &mockProjectRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Project, error) { return &model.Project{ID: id}, nil },
		}
		appRepo := &mockAppRepository{
			GetByProjectFn: func(ctx context.Context, id uint) ([]model.App, error) { return []model.App{{ID: 1}}, nil },
		}
		svc := newTestAppService(appRepo, projectRepo, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		apps, err := svc.GetAppsByProject(context.Background(), 1)
		if err != nil || len(apps) != 1 {
			t.Fatalf("unexpected result: %+v, err %v", apps, err)
		}
	})
}

func TestAppService_GetAppsByTemplate(t *testing.T) {
	t.Run("returns ErrNotFound when template missing", func(t *testing.T) {
		templateRepo := &mockTemplateRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.Template, error) { return nil, nil },
		}
		svc := newTestAppService(&mockAppRepository{}, &mockProjectRepository{}, templateRepo, &mockForge{}, &mockCI{})

		_, err := svc.GetAppsByTemplate(context.Background(), 1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestAppService_RepairWebhook(t *testing.T) {
	t.Run("returns ErrNotFound when app missing", func(t *testing.T) {
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return nil, nil },
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		err := svc.RepairWebhook(context.Background(), 1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("errors on invalid repo full name", func(t *testing.T) {
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) {
				return &model.App{ID: 1, RepoFullName: "no-slash"}, nil
			},
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		err := svc.RepairWebhook(context.Background(), 1)
		if err == nil {
			t.Fatal("expected error for invalid repo full name")
		}
	})

	t.Run("recreates webhook", func(t *testing.T) {
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) {
				return &model.App{ID: 1, RepoFullName: "owner/app"}, nil
			},
		}
		var deletedCalled, createdCalled bool
		forge := &mockForge{
			DeleteWebhookFn: func(ctx context.Context, owner, repo, callbackURL string) error {
				deletedCalled = true
				return nil
			},
			CreateWebhookFn: func(ctx context.Context, owner, repo, callbackURL, secret, branch string) error {
				createdCalled = true
				return nil
			},
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, forge, &mockCI{})

		if err := svc.RepairWebhook(context.Background(), 1); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !deletedCalled || !createdCalled {
			t.Fatalf("expected both delete and create webhook to be called")
		}
	})
}

func TestAppService_UpdateApp(t *testing.T) {
	t.Run("returns ErrNotFound when app missing", func(t *testing.T) {
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return nil, nil },
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		_, err := svc.UpdateApp(context.Background(), 1, 1, false, "name", "desc", nil)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden for non-admin non-creator", func(t *testing.T) {
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) {
				return &model.App{ID: 1, CreatorID: 42}, nil
			},
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		_, err := svc.UpdateApp(context.Background(), 1, 999, false, "name", "desc", nil)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("updates fields for creator", func(t *testing.T) {
		app := &model.App{ID: 1, CreatorID: 42, Name: "old"}
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return app, nil },
			UpdateFn:  func(ctx context.Context, app *model.App) error { return nil },
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		projectID := uint(3)
		got, err := svc.UpdateApp(context.Background(), 1, 42, false, "new", "new desc", &projectID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "new" || got.Description != "new desc" || got.ProjectID == nil || *got.ProjectID != 3 {
			t.Fatalf("app not updated as expected: %+v", got)
		}
	})
}

func TestAppService_DeleteApp(t *testing.T) {
	t.Run("returns ErrNotFound when app missing", func(t *testing.T) {
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return nil, nil },
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		err := svc.DeleteApp(context.Background(), 1, 1, false)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden for non-admin non-creator", func(t *testing.T) {
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) {
				return &model.App{ID: 1, CreatorID: 42}, nil
			},
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, &mockForge{}, &mockCI{})

		err := svc.DeleteApp(context.Background(), 1, 999, false)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("propagates ci delete failure without deleting forge repo", func(t *testing.T) {
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) {
				return &model.App{ID: 1, CreatorID: 42}, nil
			},
		}
		forgeDeleteCalled := false
		forge := &mockForge{
			DeleteRepoFn: func(ctx context.Context, owner, name string) error {
				forgeDeleteCalled = true
				return nil
			},
		}
		ci := &mockCI{
			DeleteRepoFn: func(ctx context.Context, ciRepoID int64, slug string) error { return errBoom },
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, forge, ci)

		err := svc.DeleteApp(context.Background(), 1, 42, false)
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
		if forgeDeleteCalled {
			t.Fatal("forge repo should not be deleted when ci deletion fails")
		}
	})

	t.Run("deletes ci repo, forge repo, and db record", func(t *testing.T) {
		app := &model.App{ID: 1, CreatorID: 42, Name: "app"}
		var appDeleted bool
		appRepo := &mockAppRepository{
			GetByIDFn: func(ctx context.Context, id uint) (*model.App, error) { return app, nil },
			DeleteFn: func(ctx context.Context, a *model.App) error {
				appDeleted = true
				return nil
			},
		}
		forge := &mockForge{
			DeleteRepoFn: func(ctx context.Context, owner, name string) error { return nil },
		}
		ci := &mockCI{
			DeleteRepoFn: func(ctx context.Context, ciRepoID int64, slug string) error { return nil },
		}
		svc := newTestAppService(appRepo, &mockProjectRepository{}, &mockTemplateRepository{}, forge, ci)

		if err := svc.DeleteApp(context.Background(), 1, 42, true); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !appDeleted {
			t.Fatal("expected app repository Delete to be called")
		}
	})
}

func ptr[T any](v T) *T { return &v }
