package repository

import (
	"context"
	"testing"
)

func TestServerConfigRepository_GetMissingKey(t *testing.T) {
	db := newTestDB(t)
	repo := NewServerConfigRepository(db)

	value, err := repo.Get(context.Background(), "does-not-exist")
	if err != nil {
		t.Fatalf("expected no error for a missing key, got %v", err)
	}
	if value != "" {
		t.Fatalf("expected empty value for a missing key, got %q", value)
	}
}

func TestServerConfigRepository_SetAndGet(t *testing.T) {
	db := newTestDB(t)
	repo := NewServerConfigRepository(db)

	if err := repo.Set(context.Background(), "forge_type", "gitea"); err != nil {
		t.Fatalf("failed to set config: %v", err)
	}

	value, err := repo.Get(context.Background(), "forge_type")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "gitea" {
		t.Fatalf("expected %q, got %q", "gitea", value)
	}
}

func TestServerConfigRepository_SetUpserts(t *testing.T) {
	db := newTestDB(t)
	repo := NewServerConfigRepository(db)

	if err := repo.Set(context.Background(), "forge_type", "gitea"); err != nil {
		t.Fatalf("failed to set config: %v", err)
	}
	if err := repo.Set(context.Background(), "forge_type", "github"); err != nil {
		t.Fatalf("failed to overwrite config: %v", err)
	}

	value, err := repo.Get(context.Background(), "forge_type")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "github" {
		t.Fatalf("expected the second Set to overwrite the first, got %q", value)
	}
}
