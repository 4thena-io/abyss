package repository

import (
	"context"
	"testing"

	"github.com/4thena-io/abyss/internal/model"
)

func TestUserRepository_SaveAndGetByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{ForgeID: 1, Username: "alice", Email: "alice@example.com"}
	if err := repo.Save(context.Background(), user); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}
	if user.ID == 0 {
		t.Fatal("expected user ID to be set after save")
	}

	got, err := repo.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Username != "alice" {
		t.Fatalf("expected to find saved user, got %+v", got)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)

	got, err := repo.GetByID(context.Background(), 999)
	if err != nil {
		t.Fatalf("expected no error for missing user, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for missing user, got %+v", got)
	}
}

func TestUserRepository_GetByForgeID(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{ForgeID: 42, Username: "bob"}
	if err := repo.Save(context.Background(), user); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	got, err := repo.GetByForgeID(context.Background(), 42)
	if err != nil || got == nil {
		t.Fatalf("expected to find user by forge id, got %+v, err %v", got, err)
	}

	missing, err := repo.GetByForgeID(context.Background(), 999)
	if err != nil || missing != nil {
		t.Fatalf("expected nil for missing forge id, got %+v, err %v", missing, err)
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{ForgeID: 1, Username: "carol"}
	if err := repo.Save(context.Background(), user); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	got, err := repo.GetByUsername(context.Background(), "carol")
	if err != nil || got == nil {
		t.Fatalf("expected to find user by username, got %+v, err %v", got, err)
	}

	missing, err := repo.GetByUsername(context.Background(), "does-not-exist")
	if err != nil || missing != nil {
		t.Fatalf("expected nil for missing username, got %+v, err %v", missing, err)
	}
}

func TestUserRepository_GetByToken(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)

	token := "pat-abc123"
	user := &model.User{ForgeID: 1, Username: "dave", Token: &token}
	if err := repo.Save(context.Background(), user); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	got, err := repo.GetByToken(context.Background(), token)
	if err != nil || got == nil {
		t.Fatalf("expected to find user by token, got %+v, err %v", got, err)
	}

	missing, err := repo.GetByToken(context.Background(), "does-not-exist")
	if err != nil || missing != nil {
		t.Fatalf("expected nil for missing token, got %+v, err %v", missing, err)
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{ForgeID: 1, Username: "erin", Email: "old@example.com"}
	if err := repo.Save(context.Background(), user); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	user.Email = "new@example.com"
	if err := repo.Update(context.Background(), user); err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	got, err := repo.GetByID(context.Background(), user.ID)
	if err != nil || got.Email != "new@example.com" {
		t.Fatalf("expected updated email, got %+v, err %v", got, err)
	}
}

func TestUserRepository_GetAllAndCount(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)

	for i := 0; i < 3; i++ {
		user := &model.User{ForgeID: int64(i + 1), Username: "user" + string(rune('a'+i))}
		if err := repo.Save(context.Background(), user); err != nil {
			t.Fatalf("failed to save user: %v", err)
		}
	}

	users, err := repo.GetAll(context.Background())
	if err != nil || len(users) != 3 {
		t.Fatalf("expected 3 users, got %d, err %v", len(users), err)
	}

	count, err := repo.Count(context.Background())
	if err != nil || count != 3 {
		t.Fatalf("expected count 3, got %d, err %v", count, err)
	}
}
