package repository

import (
	"context"
	"testing"

	"github.com/4thena-io/abyss/internal/model"
)

func TestProjectRepository_SaveAndGetByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewProjectRepository(db)
	user := createTestUser(t, db, "creator")

	project := &model.Project{Name: "my-project", Description: "desc", CreatorID: user.ID}
	if err := repo.Save(context.Background(), project); err != nil {
		t.Fatalf("failed to save project: %v", err)
	}
	if project.ID == 0 {
		t.Fatal("expected project ID to be set after save")
	}

	got, err := repo.GetByID(context.Background(), project.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Name != "my-project" {
		t.Fatalf("expected to find saved project, got %+v", got)
	}
	if got.Creator.Username != "creator" {
		t.Fatalf("expected creator to be preloaded, got %+v", got.Creator)
	}
}

func TestProjectRepository_GetByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewProjectRepository(db)

	got, err := repo.GetByID(context.Background(), 999)
	if err != nil {
		t.Fatalf("expected no error for missing project, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for missing project, got %+v", got)
	}
}

func TestProjectRepository_GetByName(t *testing.T) {
	db := newTestDB(t)
	repo := NewProjectRepository(db)
	user := createTestUser(t, db, "creator")

	project := &model.Project{Name: "unique-project", CreatorID: user.ID}
	if err := repo.Save(context.Background(), project); err != nil {
		t.Fatalf("failed to save project: %v", err)
	}

	got, err := repo.GetByName(context.Background(), "unique-project")
	if err != nil || got == nil {
		t.Fatalf("expected to find project by name, got %+v, err %v", got, err)
	}

	missing, err := repo.GetByName(context.Background(), "does-not-exist")
	if err != nil || missing != nil {
		t.Fatalf("expected nil for missing name, got %+v, err %v", missing, err)
	}
}

func TestProjectRepository_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewProjectRepository(db)
	user := createTestUser(t, db, "creator")

	project := &model.Project{Name: "project", CreatorID: user.ID}
	if err := repo.Save(context.Background(), project); err != nil {
		t.Fatalf("failed to save project: %v", err)
	}

	project.Description = "updated"
	if err := repo.Update(context.Background(), project); err != nil {
		t.Fatalf("failed to update project: %v", err)
	}

	got, err := repo.GetByID(context.Background(), project.ID)
	if err != nil || got.Description != "updated" {
		t.Fatalf("expected updated description, got %+v, err %v", got, err)
	}
}

func TestProjectRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewProjectRepository(db)
	user := createTestUser(t, db, "creator")

	project := &model.Project{Name: "project", CreatorID: user.ID}
	if err := repo.Save(context.Background(), project); err != nil {
		t.Fatalf("failed to save project: %v", err)
	}

	if err := repo.Delete(context.Background(), project); err != nil {
		t.Fatalf("failed to delete project: %v", err)
	}

	got, err := repo.GetByID(context.Background(), project.ID)
	if err != nil || got != nil {
		t.Fatalf("expected project to be gone after delete, got %+v, err %v", got, err)
	}
}

func TestProjectRepository_GetByTeamAndCountByTeam(t *testing.T) {
	db := newTestDB(t)
	repo := NewProjectRepository(db)
	user := createTestUser(t, db, "creator")

	team := &model.Team{Name: "team"}
	if err := db.Create(team).Error; err != nil {
		t.Fatalf("failed to create team: %v", err)
	}

	for i := 0; i < 2; i++ {
		project := &model.Project{Name: "project" + string(rune('a'+i)), CreatorID: user.ID, TeamID: &team.ID}
		if err := repo.Save(context.Background(), project); err != nil {
			t.Fatalf("failed to save project: %v", err)
		}
	}

	projects, err := repo.GetByTeam(context.Background(), team.ID)
	if err != nil || len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d, err %v", len(projects), err)
	}

	count, err := repo.CountByTeam(context.Background(), team.ID)
	if err != nil || count != 2 {
		t.Fatalf("expected count 2, got %d, err %v", count, err)
	}
}

func TestProjectRepository_GetAll(t *testing.T) {
	db := newTestDB(t)
	repo := NewProjectRepository(db)
	user := createTestUser(t, db, "creator")

	for i := 0; i < 3; i++ {
		project := &model.Project{Name: "project" + string(rune('a'+i)), CreatorID: user.ID}
		if err := repo.Save(context.Background(), project); err != nil {
			t.Fatalf("failed to save project: %v", err)
		}
	}

	projects, err := repo.GetAll(context.Background())
	if err != nil || len(projects) != 3 {
		t.Fatalf("expected 3 projects, got %d, err %v", len(projects), err)
	}
}
