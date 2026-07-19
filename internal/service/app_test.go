package service

import (
	"context"
	"errors"
	"testing"

	cimocks "github.com/4thena-io/abyss/internal/integration/ci/mocks"
	forgemocks "github.com/4thena-io/abyss/internal/integration/forge/mocks"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service/mocks"
	"github.com/stretchr/testify/mock"
)

func newTestAppService(appRepo AppRepository, projectRepo ProjectRepository, templateRepo TemplateRepository, forge *forgemocks.MockForge, ci *cimocks.MockCI) *AppService {
	return NewAppService(appRepo, projectRepo, templateRepo, forge, ci, nil, "owner", "https://abyss.example", "secret", "main")
}

func TestAppService_CreateAppFromTemplate(t *testing.T) {
	t.Run("errors when project id set but project missing", func(t *testing.T) {
		projectID := uint(5)
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByID(mock.Anything, uint(5)).Return(nil, nil)
		svc := newTestAppService(mocks.NewMockAppRepository(t), projectRepo, mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		_, err := svc.CreateAppFromTemplate(context.Background(), &model.App{Name: "app", ProjectID: &projectID, TemplateID: ptr(uint(1))})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("errors when template id is nil", func(t *testing.T) {
		svc := newTestAppService(mocks.NewMockAppRepository(t), mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		_, err := svc.CreateAppFromTemplate(context.Background(), &model.App{Name: "app"})
		if err == nil {
			t.Fatal("expected error for missing template ID")
		}
	})

	t.Run("errors when template does not exist", func(t *testing.T) {
		templateRepo := mocks.NewMockTemplateRepository(t)
		templateRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		svc := newTestAppService(mocks.NewMockAppRepository(t), mocks.NewMockProjectRepository(t), templateRepo, forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		_, err := svc.CreateAppFromTemplate(context.Background(), &model.App{Name: "app", TemplateID: ptr(uint(1))})
		if err == nil {
			t.Fatal("expected error for missing template")
		}
	})

	t.Run("returns ErrConflict when app name already exists", func(t *testing.T) {
		templateRepo := mocks.NewMockTemplateRepository(t)
		templateRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.Template{ID: 1, CloneURL: "https://forge/template.git"}, nil)
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByName(mock.Anything, "taken").Return(&model.App{ID: 1, Name: "taken"}, nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), templateRepo, forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		_, err := svc.CreateAppFromTemplate(context.Background(), &model.App{Name: "taken", TemplateID: ptr(uint(1))})
		if !errors.Is(err, ErrConflict) {
			t.Fatalf("expected ErrConflict, got %v", err)
		}
	})

	t.Run("propagates forge repo creation failure without touching git", func(t *testing.T) {
		templateRepo := mocks.NewMockTemplateRepository(t)
		templateRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.Template{ID: 1, CloneURL: "https://forge/template.git"}, nil)
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByName(mock.Anything, "app").Return(nil, nil)
		forge := forgemocks.NewMockForge(t)
		forge.EXPECT().CreateRepo(mock.Anything, "owner", "app").Return(nil, errBoom)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), templateRepo, forge, cimocks.NewMockCI(t))

		_, err := svc.CreateAppFromTemplate(context.Background(), &model.App{Name: "app", TemplateID: ptr(uint(1))})
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})
}

func TestAppService_CreateAppFromRepo(t *testing.T) {
	t.Run("errors when project id set but project missing", func(t *testing.T) {
		projectID := uint(5)
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByID(mock.Anything, uint(5)).Return(nil, nil)
		svc := newTestAppService(mocks.NewMockAppRepository(t), projectRepo, mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		_, err := svc.CreateAppFromRepo(context.Background(), &model.App{Name: "app", ProjectID: &projectID})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns ErrConflict when app name already exists", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByName(mock.Anything, "taken").Return(&model.App{ID: 1, Name: "taken"}, nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		_, err := svc.CreateAppFromRepo(context.Background(), &model.App{Name: "taken", RepoID: 99})
		if !errors.Is(err, ErrConflict) {
			t.Fatalf("expected ErrConflict, got %v", err)
		}
	})

	t.Run("errors when repo does not exist in forge", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByName(mock.Anything, "app").Return(nil, nil)
		forge := forgemocks.NewMockForge(t)
		forge.EXPECT().GetRepo(mock.Anything, int64(99)).Return(nil, errBoom)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forge, cimocks.NewMockCI(t))

		_, err := svc.CreateAppFromRepo(context.Background(), &model.App{Name: "app", RepoID: 99})
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})
}

func TestAppService_GetAppByID(t *testing.T) {
	appRepo := mocks.NewMockAppRepository(t)
	appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.App{ID: 1}, nil)
	svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

	got, err := svc.GetAppByID(context.Background(), 1)
	if err != nil || got == nil {
		t.Fatalf("expected app, got %+v, err %v", got, err)
	}
}

func TestAppService_GetAppBuilds(t *testing.T) {
	t.Run("returns ErrNotFound when app missing", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		_, err := svc.GetAppBuilds(context.Background(), 1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns builds from ci", func(t *testing.T) {
		app := &model.App{ID: 1, CIID: 55, CISlug: "owner/app"}
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(app, nil)
		ci := cimocks.NewMockCI(t)
		ci.EXPECT().GetBuilds(mock.Anything, int64(55), "owner/app").Return([]model.Build{{ID: 1}}, nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), ci)

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
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		svc := newTestAppService(mocks.NewMockAppRepository(t), projectRepo, mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		_, err := svc.GetAppsByProject(context.Background(), 1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns apps for project", func(t *testing.T) {
		projectRepo := mocks.NewMockProjectRepository(t)
		projectRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.Project{ID: 1}, nil)
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByProject(mock.Anything, uint(1)).Return([]model.App{{ID: 1}}, nil)
		svc := newTestAppService(appRepo, projectRepo, mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		apps, err := svc.GetAppsByProject(context.Background(), 1)
		if err != nil || len(apps) != 1 {
			t.Fatalf("unexpected result: %+v, err %v", apps, err)
		}
	})
}

func TestAppService_GetAppsByTemplate(t *testing.T) {
	t.Run("returns ErrNotFound when template missing", func(t *testing.T) {
		templateRepo := mocks.NewMockTemplateRepository(t)
		templateRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		svc := newTestAppService(mocks.NewMockAppRepository(t), mocks.NewMockProjectRepository(t), templateRepo, forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		_, err := svc.GetAppsByTemplate(context.Background(), 1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestAppService_RepairWebhook(t *testing.T) {
	t.Run("returns ErrNotFound when app missing", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		err := svc.RepairWebhook(context.Background(), 1)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("errors on invalid repo full name", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.App{ID: 1, RepoFullName: "no-slash"}, nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		err := svc.RepairWebhook(context.Background(), 1)
		if err == nil {
			t.Fatal("expected error for invalid repo full name")
		}
	})

	t.Run("recreates webhook", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.App{ID: 1, RepoFullName: "owner/app"}, nil)
		forge := forgemocks.NewMockForge(t)
		forge.EXPECT().DeleteWebhook(mock.Anything, "owner", "app", mock.Anything).Return(nil)
		forge.EXPECT().CreateWebhook(mock.Anything, "owner", "app", mock.Anything, "secret", "main").Return(nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forge, cimocks.NewMockCI(t))

		if err := svc.RepairWebhook(context.Background(), 1); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestAppService_UpdateApp(t *testing.T) {
	t.Run("returns ErrNotFound when app missing", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		_, err := svc.UpdateApp(context.Background(), 1, 1, false, "name", "desc", nil)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden for non-admin non-creator", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.App{ID: 1, CreatorID: 42}, nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		_, err := svc.UpdateApp(context.Background(), 1, 999, false, "name", "desc", nil)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("updates fields for creator", func(t *testing.T) {
		app := &model.App{ID: 1, CreatorID: 42, Name: "old"}
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(app, nil)
		appRepo.EXPECT().Update(mock.Anything, app).Return(nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

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
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(nil, nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		err := svc.DeleteApp(context.Background(), 1, 1, false)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrForbidden for non-admin non-creator", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.App{ID: 1, CreatorID: 42}, nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forgemocks.NewMockForge(t), cimocks.NewMockCI(t))

		err := svc.DeleteApp(context.Background(), 1, 999, false)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("propagates ci delete failure without deleting forge repo", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(&model.App{ID: 1, CreatorID: 42}, nil)
		forge := forgemocks.NewMockForge(t)
		ci := cimocks.NewMockCI(t)
		ci.EXPECT().DeleteRepo(mock.Anything, mock.Anything, mock.Anything).Return(errBoom)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forge, ci)

		err := svc.DeleteApp(context.Background(), 1, 42, false)
		if !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
		// forge.DeleteRepo has no expectation set, so the mock will panic via
		// its testing.T handle if the service calls it despite the ci
		// deletion failure — no explicit assertion needed beyond that.
	})

	t.Run("deletes ci repo, forge repo, and db record", func(t *testing.T) {
		app := &model.App{ID: 1, CreatorID: 42, Name: "app"}
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(mock.Anything, uint(1)).Return(app, nil)
		appRepo.EXPECT().Delete(mock.Anything, app).Return(nil)
		forge := forgemocks.NewMockForge(t)
		forge.EXPECT().DeleteRepo(mock.Anything, "owner", "app").Return(nil)
		ci := cimocks.NewMockCI(t)
		ci.EXPECT().DeleteRepo(mock.Anything, mock.Anything, mock.Anything).Return(nil)
		svc := newTestAppService(appRepo, mocks.NewMockProjectRepository(t), mocks.NewMockTemplateRepository(t), forge, ci)

		if err := svc.DeleteApp(context.Background(), 1, 42, true); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
