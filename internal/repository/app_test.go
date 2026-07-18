package repository

import (
	"context"
	"testing"

	"github.com/4thena-io/abyss/internal/model"
	"gorm.io/gorm"
)

func createTestUser(t *testing.T, db *gorm.DB, username string) model.User {
	t.Helper()
	user := model.User{ForgeID: int64(len(username)), Username: username}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return user
}

func TestAppRepository_SaveAndGetByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewAppRepository(db)
	user := createTestUser(t, db, "creator")

	app := &model.App{
		Name:         "my-app",
		Description:  "desc",
		RepoID:       1,
		RepoFullName: "owner/my-app",
		RepoURL:      "https://forge/owner/my-app",
		CloneURL:     "https://forge/owner/my-app.git",
		CIID:         1,
		CIURL:        "https://ci/owner/my-app",
		CreatorID:    user.ID,
	}
	if err := repo.Save(context.Background(), app); err != nil {
		t.Fatalf("failed to save app: %v", err)
	}
	if app.ID == 0 {
		t.Fatal("expected app ID to be set after save")
	}

	got, err := repo.GetByID(context.Background(), app.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Name != "my-app" {
		t.Fatalf("expected to find saved app, got %+v", got)
	}
	if got.Creator.Username != "creator" {
		t.Fatalf("expected creator to be preloaded, got %+v", got.Creator)
	}
}

func TestAppRepository_GetByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewAppRepository(db)

	got, err := repo.GetByID(context.Background(), 999)
	if err != nil {
		t.Fatalf("expected no error for missing app, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for missing app, got %+v", got)
	}
}

func TestAppRepository_GetByName(t *testing.T) {
	db := newTestDB(t)
	repo := NewAppRepository(db)
	user := createTestUser(t, db, "creator")

	app := &model.App{Name: "unique-app", RepoID: 1, RepoFullName: "o/a", RepoURL: "u", CloneURL: "c", CIID: 1, CIURL: "u", CreatorID: user.ID}
	if err := repo.Save(context.Background(), app); err != nil {
		t.Fatalf("failed to save app: %v", err)
	}

	got, err := repo.GetByName(context.Background(), "unique-app")
	if err != nil || got == nil {
		t.Fatalf("expected to find app by name, got %+v, err %v", got, err)
	}

	missing, err := repo.GetByName(context.Background(), "does-not-exist")
	if err != nil || missing != nil {
		t.Fatalf("expected nil for missing name, got %+v, err %v", missing, err)
	}
}

func TestAppRepository_GetByProjectAndTemplate(t *testing.T) {
	db := newTestDB(t)
	appRepo := NewAppRepository(db)
	user := createTestUser(t, db, "creator")

	project := &model.Project{Name: "proj", CreatorID: user.ID}
	if err := db.Create(project).Error; err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	template := &model.Template{Name: "tmpl", CreatorID: user.ID}
	if err := db.Create(template).Error; err != nil {
		t.Fatalf("failed to create template: %v", err)
	}

	app := &model.App{
		Name: "scoped-app", RepoID: 1, RepoFullName: "o/a", RepoURL: "u", CloneURL: "c", CIID: 1, CIURL: "u",
		CreatorID: user.ID, ProjectID: &project.ID, TemplateID: &template.ID,
	}
	if err := appRepo.Save(context.Background(), app); err != nil {
		t.Fatalf("failed to save app: %v", err)
	}

	byProject, err := appRepo.GetByProject(context.Background(), project.ID)
	if err != nil || len(byProject) != 1 {
		t.Fatalf("expected 1 app by project, got %d, err %v", len(byProject), err)
	}

	byTemplate, err := appRepo.GetByTemplate(context.Background(), template.ID)
	if err != nil || len(byTemplate) != 1 {
		t.Fatalf("expected 1 app by template, got %d, err %v", len(byTemplate), err)
	}
}

func TestAppRepository_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewAppRepository(db)
	user := createTestUser(t, db, "creator")

	app := &model.App{Name: "app", RepoID: 1, RepoFullName: "o/a", RepoURL: "u", CloneURL: "c", CIID: 1, CIURL: "u", CreatorID: user.ID}
	if err := repo.Save(context.Background(), app); err != nil {
		t.Fatalf("failed to save app: %v", err)
	}

	app.Description = "updated description"
	if err := repo.Update(context.Background(), app); err != nil {
		t.Fatalf("failed to update app: %v", err)
	}

	got, err := repo.GetByID(context.Background(), app.ID)
	if err != nil || got.Description != "updated description" {
		t.Fatalf("expected updated description, got %+v, err %v", got, err)
	}
}

func TestAppRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewAppRepository(db)
	user := createTestUser(t, db, "creator")

	app := &model.App{Name: "app", RepoID: 1, RepoFullName: "o/a", RepoURL: "u", CloneURL: "c", CIID: 1, CIURL: "u", CreatorID: user.ID}
	if err := repo.Save(context.Background(), app); err != nil {
		t.Fatalf("failed to save app: %v", err)
	}

	if err := repo.Delete(context.Background(), app); err != nil {
		t.Fatalf("failed to delete app: %v", err)
	}

	got, err := repo.GetByID(context.Background(), app.ID)
	if err != nil || got != nil {
		t.Fatalf("expected app to be gone after delete, got %+v, err %v", got, err)
	}
}

func TestAppRepository_CountByTeam(t *testing.T) {
	db := newTestDB(t)
	appRepo := NewAppRepository(db)
	user := createTestUser(t, db, "creator")

	team := &model.Team{Name: "team"}
	if err := db.Create(team).Error; err != nil {
		t.Fatalf("failed to create team: %v", err)
	}
	project := &model.Project{Name: "proj", CreatorID: user.ID, TeamID: &team.ID}
	if err := db.Create(project).Error; err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	for i := 0; i < 3; i++ {
		app := &model.App{
			Name: "app" + string(rune('a'+i)), RepoID: int64(i + 1), RepoFullName: "o/a", RepoURL: "u", CloneURL: "c", CIID: 1, CIURL: "u",
			CreatorID: user.ID, ProjectID: &project.ID,
		}
		if err := appRepo.Save(context.Background(), app); err != nil {
			t.Fatalf("failed to save app: %v", err)
		}
	}

	count, err := appRepo.CountByTeam(context.Background(), team.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected count 3, got %d", count)
	}
}
