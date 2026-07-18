package repository

import (
	"context"
	"testing"
	"time"

	"github.com/4thena-io/abyss/internal/model"
	"gorm.io/gorm"
)

func createTestApp(t *testing.T, db *gorm.DB, name string, creatorID uint) model.App {
	t.Helper()
	app := model.App{
		Name: name, RepoID: int64(len(name)), RepoFullName: "o/" + name, RepoURL: "u", CloneURL: "c",
		CIID: 1, CIURL: "u", CreatorID: creatorID,
	}
	if err := db.Create(&app).Error; err != nil {
		t.Fatalf("failed to create test app: %v", err)
	}
	return app
}

func TestDeploymentRepository_SaveAndGetByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewDeploymentRepository(db)
	user := createTestUser(t, db, "creator")
	app := createTestApp(t, db, "app", user.ID)

	deployment := &model.Deployment{
		AppID: app.ID, Environment: "production", Status: "success", Commit: "abc123",
		TriggeredBy: "alice", Duration: 42, DeployedAt: time.Now(),
	}
	if err := repo.Save(context.Background(), deployment); err != nil {
		t.Fatalf("failed to save deployment: %v", err)
	}
	if deployment.ID == 0 {
		t.Fatal("expected deployment ID to be set after save")
	}

	got, err := repo.GetByID(context.Background(), deployment.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Environment != "production" {
		t.Fatalf("expected to find saved deployment, got %+v", got)
	}
}

func TestDeploymentRepository_GetByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewDeploymentRepository(db)

	got, err := repo.GetByID(context.Background(), 999)
	if err != nil {
		t.Fatalf("expected no error for missing deployment, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for missing deployment, got %+v", got)
	}
}

func TestDeploymentRepository_GetByApp_OrderedByDeployedAtDesc(t *testing.T) {
	db := newTestDB(t)
	repo := NewDeploymentRepository(db)
	user := createTestUser(t, db, "creator")
	app := createTestApp(t, db, "app", user.ID)

	older := &model.Deployment{AppID: app.ID, Environment: "prod", Status: "success", Commit: "a", TriggeredBy: "x", DeployedAt: time.Now().Add(-time.Hour)}
	newer := &model.Deployment{AppID: app.ID, Environment: "prod", Status: "success", Commit: "b", TriggeredBy: "x", DeployedAt: time.Now()}
	if err := repo.Save(context.Background(), older); err != nil {
		t.Fatalf("failed to save older deployment: %v", err)
	}
	if err := repo.Save(context.Background(), newer); err != nil {
		t.Fatalf("failed to save newer deployment: %v", err)
	}

	got, err := repo.GetByApp(context.Background(), app.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].Commit != "b" || got[1].Commit != "a" {
		t.Fatalf("expected deployments ordered newest first, got %+v", got)
	}
}

func TestDeploymentRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewDeploymentRepository(db)
	user := createTestUser(t, db, "creator")
	app := createTestApp(t, db, "app", user.ID)

	deployment := &model.Deployment{AppID: app.ID, Environment: "prod", Status: "success", Commit: "a", TriggeredBy: "x", DeployedAt: time.Now()}
	if err := repo.Save(context.Background(), deployment); err != nil {
		t.Fatalf("failed to save deployment: %v", err)
	}

	if err := repo.Delete(context.Background(), deployment); err != nil {
		t.Fatalf("failed to delete deployment: %v", err)
	}

	got, err := repo.GetByID(context.Background(), deployment.ID)
	if err != nil || got != nil {
		t.Fatalf("expected deployment to be gone after delete, got %+v, err %v", got, err)
	}
}
