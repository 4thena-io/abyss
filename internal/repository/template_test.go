package repository

import (
	"context"
	"testing"

	"github.com/4thena-io/abyss/internal/model"
)

func TestTemplateRepository_SaveAndGetByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewTemplateRepository(db)
	user := createTestUser(t, db, "creator")

	template := &model.Template{Name: "go-service", Description: "desc", Kind: "service", Language: "go", CreatorID: user.ID}
	if err := repo.Save(context.Background(), template); err != nil {
		t.Fatalf("failed to save template: %v", err)
	}
	if template.ID == 0 {
		t.Fatal("expected template ID to be set after save")
	}

	got, err := repo.GetByID(context.Background(), template.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Name != "go-service" {
		t.Fatalf("expected to find saved template, got %+v", got)
	}
	if got.Creator.Username != "creator" {
		t.Fatalf("expected creator to be preloaded, got %+v", got.Creator)
	}
}

func TestTemplateRepository_GetByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewTemplateRepository(db)

	got, err := repo.GetByID(context.Background(), 999)
	if err != nil {
		t.Fatalf("expected no error for missing template, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for missing template, got %+v", got)
	}
}

func TestTemplateRepository_GetByName(t *testing.T) {
	db := newTestDB(t)
	repo := NewTemplateRepository(db)
	user := createTestUser(t, db, "creator")

	template := &model.Template{Name: "unique-template", CreatorID: user.ID}
	if err := repo.Save(context.Background(), template); err != nil {
		t.Fatalf("failed to save template: %v", err)
	}

	got, err := repo.GetByName(context.Background(), "unique-template")
	if err != nil || got == nil {
		t.Fatalf("expected to find template by name, got %+v, err %v", got, err)
	}

	missing, err := repo.GetByName(context.Background(), "does-not-exist")
	if err != nil || missing != nil {
		t.Fatalf("expected nil for missing name, got %+v, err %v", missing, err)
	}
}

func TestTemplateRepository_GetAll(t *testing.T) {
	db := newTestDB(t)
	repo := NewTemplateRepository(db)
	user := createTestUser(t, db, "creator")

	for i := 0; i < 3; i++ {
		template := &model.Template{Name: "template" + string(rune('a'+i)), CreatorID: user.ID}
		if err := repo.Save(context.Background(), template); err != nil {
			t.Fatalf("failed to save template: %v", err)
		}
	}

	templates, err := repo.GetAll(context.Background())
	if err != nil || len(templates) != 3 {
		t.Fatalf("expected 3 templates, got %d, err %v", len(templates), err)
	}
}

func TestTemplateRepository_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewTemplateRepository(db)
	user := createTestUser(t, db, "creator")

	template := &model.Template{Name: "template", CreatorID: user.ID}
	if err := repo.Save(context.Background(), template); err != nil {
		t.Fatalf("failed to save template: %v", err)
	}

	template.Description = "updated"
	if err := repo.Update(context.Background(), template); err != nil {
		t.Fatalf("failed to update template: %v", err)
	}

	got, err := repo.GetByID(context.Background(), template.ID)
	if err != nil || got.Description != "updated" {
		t.Fatalf("expected updated description, got %+v, err %v", got, err)
	}
}

func TestTemplateRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewTemplateRepository(db)
	user := createTestUser(t, db, "creator")

	template := &model.Template{Name: "template", CreatorID: user.ID}
	if err := repo.Save(context.Background(), template); err != nil {
		t.Fatalf("failed to save template: %v", err)
	}

	if err := repo.Delete(context.Background(), template); err != nil {
		t.Fatalf("failed to delete template: %v", err)
	}

	got, err := repo.GetByID(context.Background(), template.ID)
	if err != nil || got != nil {
		t.Fatalf("expected template to be gone after delete, got %+v, err %v", got, err)
	}
}
