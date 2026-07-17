package git

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type GitClient struct {
	token  string
	name   string
	email  string
	branch string
}

func New(token, name, email, branch string) *GitClient {
	if branch == "" {
		branch = "main"
	}
	return &GitClient{token: token, name: name, email: email, branch: branch}
}

// Clone clones a repository to a local path on the configured branch.
func (c *GitClient) Clone(repoURL, localPath string) error {
	_, err := gogit.PlainClone(localPath, false, &gogit.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(c.branch),
		Auth:          c.auth(),
	})
	return err
}

func (c *GitClient) CloneSubset(repoURL, localPath string, _ []string) error {
	_, err := gogit.PlainClone(localPath, false, &gogit.CloneOptions{
		URL:   repoURL,
		Depth: 1,
		Auth:  c.auth(),
	})
	return err
}

// InitAndPush initializes a new git repo, adds all files, commits, and pushes to remote.
func (c *GitClient) InitAndPush(localPath, remoteURL, commitMsg string) error {
	repo, err := gogit.PlainInitWithOptions(localPath, &gogit.PlainInitOptions{
		InitOptions: gogit.InitOptions{DefaultBranch: plumbing.NewBranchReferenceName(c.branch)},
	})
	if err != nil {
		return fmt.Errorf("failed to init repo: %w", err)
	}

	// Add remote
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteURL},
	})
	if err != nil {
		return fmt.Errorf("failed to add remote: %w", err)
	}

	// Add all files
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.AddGlob(".")
	if err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	// Commit
	_, err = worktree.Commit(commitMsg, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  c.name,
			Email: c.email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	// Push to the configured branch.
	refspec := config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", c.branch, c.branch))
	err = repo.Push(&gogit.PushOptions{
		Auth:     c.auth(),
		RefSpecs: []config.RefSpec{refspec},
	})
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

// FileExists checks if a file exists in the local path.
func (c *GitClient) FileExists(localPath, filePath string) bool {
	fullPath := filepath.Join(localPath, filePath)
	_, err := os.Stat(fullPath)
	return err == nil
}

// CreateFile creates a file with the given content.
func (c *GitClient) CreateFile(localPath, filePath, content string) error {
	fullPath := filepath.Join(localPath, filePath)

	// Ensure parent directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.WriteFile(fullPath, []byte(content), 0644)
}

// RemoveGitDir removes the .git directory from a local path.
func (c *GitClient) RemoveGitDir(localPath string) error {
	gitDir := filepath.Join(localPath, ".git")
	return os.RemoveAll(gitDir)
}

// AddCommitPush adds all changes, commits, and pushes to the existing remote.
func (c *GitClient) AddCommitPush(localPath, commitMsg string) error {
	repo, err := gogit.PlainOpen(localPath)
	if err != nil {
		return fmt.Errorf("failed to open repo: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.AddGlob(".")
	if err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	_, err = worktree.Commit(commitMsg, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  c.name,
			Email: c.email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	refspec := config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", c.branch, c.branch))
	err = repo.Push(&gogit.PushOptions{
		Auth:     c.auth(),
		RefSpecs: []config.RefSpec{refspec},
	})
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

func (c *GitClient) auth() *http.BasicAuth {
	if c.token == "" {
		return nil
	}
	return &http.BasicAuth{
		Username: "git", // Can be anything for token auth
		Password: c.token,
	}
}
