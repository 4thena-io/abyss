package repository

import (
	"context"
	"testing"

	"github.com/4thena-io/abyss/internal/model"
)

func TestTeamRepository_SaveAndGetByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewTeamRepository(db)

	team := &model.Team{Name: "platform", Description: "platform team"}
	if err := repo.Save(context.Background(), team); err != nil {
		t.Fatalf("failed to save team: %v", err)
	}
	if team.ID == 0 {
		t.Fatal("expected team ID to be set after save")
	}

	got, err := repo.GetByID(context.Background(), team.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Name != "platform" {
		t.Fatalf("expected to find saved team, got %+v", got)
	}
}

func TestTeamRepository_GetByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewTeamRepository(db)

	got, err := repo.GetByID(context.Background(), 999)
	if err != nil {
		t.Fatalf("expected no error for missing team, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for missing team, got %+v", got)
	}
}

func TestTeamRepository_GetByName(t *testing.T) {
	db := newTestDB(t)
	repo := NewTeamRepository(db)

	team := &model.Team{Name: "unique-team"}
	if err := repo.Save(context.Background(), team); err != nil {
		t.Fatalf("failed to save team: %v", err)
	}

	got, err := repo.GetByName(context.Background(), "unique-team")
	if err != nil || got == nil {
		t.Fatalf("expected to find team by name, got %+v, err %v", got, err)
	}

	missing, err := repo.GetByName(context.Background(), "does-not-exist")
	if err != nil || missing != nil {
		t.Fatalf("expected nil for missing name, got %+v, err %v", missing, err)
	}
}

func TestTeamRepository_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewTeamRepository(db)

	team := &model.Team{Name: "team"}
	if err := repo.Save(context.Background(), team); err != nil {
		t.Fatalf("failed to save team: %v", err)
	}

	team.Description = "updated"
	if err := repo.Update(context.Background(), team); err != nil {
		t.Fatalf("failed to update team: %v", err)
	}

	got, err := repo.GetByID(context.Background(), team.ID)
	if err != nil || got.Description != "updated" {
		t.Fatalf("expected updated description, got %+v, err %v", got, err)
	}
}

func TestTeamRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewTeamRepository(db)

	team := &model.Team{Name: "team"}
	if err := repo.Save(context.Background(), team); err != nil {
		t.Fatalf("failed to save team: %v", err)
	}

	if err := repo.Delete(context.Background(), team); err != nil {
		t.Fatalf("failed to delete team: %v", err)
	}

	got, err := repo.GetByID(context.Background(), team.ID)
	if err != nil || got != nil {
		t.Fatalf("expected team to be gone after delete, got %+v, err %v", got, err)
	}
}

func TestTeamRepository_Members(t *testing.T) {
	db := newTestDB(t)
	repo := NewTeamRepository(db)
	user := createTestUser(t, db, "alice")

	team := &model.Team{Name: "team"}
	if err := repo.Save(context.Background(), team); err != nil {
		t.Fatalf("failed to save team: %v", err)
	}

	member := &model.TeamMember{TeamID: team.ID, UserID: user.ID, Role: "owner"}
	if err := repo.SaveMember(context.Background(), member); err != nil {
		t.Fatalf("failed to save member: %v", err)
	}

	members, err := repo.GetMembers(context.Background(), team.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 1 || members[0].User.Username != "alice" {
		t.Fatalf("expected 1 member with preloaded user, got %+v", members)
	}

	count, err := repo.CountMembers(context.Background(), team.ID)
	if err != nil || count != 1 {
		t.Fatalf("expected count 1, got %d, err %v", count, err)
	}

	byUser, err := repo.GetByUserID(context.Background(), user.ID)
	if err != nil || len(byUser) != 1 || byUser[0].TeamID != team.ID {
		t.Fatalf("expected 1 membership for user, got %+v, err %v", byUser, err)
	}

	if err := repo.DeleteMember(context.Background(), team.ID, member.ID); err != nil {
		t.Fatalf("failed to delete member: %v", err)
	}
	count, err = repo.CountMembers(context.Background(), team.ID)
	if err != nil || count != 0 {
		t.Fatalf("expected count 0 after delete, got %d, err %v", count, err)
	}
}
