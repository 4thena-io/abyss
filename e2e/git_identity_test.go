package e2e

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/4thena-io/abyss/internal/git"
	gogit "github.com/go-git/go-git/v5"
)

func TestGitCommitIdentityRecovers_E2E(t *testing.T) {
	remote := newBareRepoWithFiles(t, nil)
	local := t.TempDir()
	online := false
	unavailable := errors.New("bot unavailable")
	client := git.NewWithIdentity("", "main", func() (string, string, error) {
		if !online {
			return "", "", unavailable
		}
		return "Abyss Bot", "bot@example.com", nil
	})
	if err := client.CreateFile(local, "README.md", "hello\n"); err != nil {
		t.Fatal(err)
	}
	if err := client.InitAndPush(local, remote, "initialize"); !errors.Is(err, unavailable) {
		t.Fatalf("expected bot outage, got %v", err)
	}
	online = true
	if err := client.AddCommitPush(local, "initialize after recovery"); err != nil {
		t.Fatal(err)
	}
	clone := filepath.Join(t.TempDir(), "clone")
	if err := client.Clone(remote, clone); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(clone, "README.md"))
	if err != nil || string(data) != "hello\n" {
		t.Fatalf("content=%q err=%v", data, err)
	}
	repo, err := gogit.PlainOpen(clone)
	if err != nil {
		t.Fatal(err)
	}
	head, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		t.Fatal(err)
	}
	if commit.Author.Name != "Abyss Bot" || commit.Author.Email != "bot@example.com" {
		t.Fatalf("unexpected author: %+v", commit.Author)
	}
}
