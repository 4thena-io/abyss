package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/4thena-io/abyss/internal/git"
	cimocks "github.com/4thena-io/abyss/internal/integration/ci/mocks"
	forgemocks "github.com/4thena-io/abyss/internal/integration/forge/mocks"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service/mocks"
)

// newBareRepoWithFiles creates a local bare git repo (no network) seeded with
// an initial commit on "main" containing the given files, and returns its
// filesystem path — usable directly as a CloneURL since go-git treats plain
// local paths as a valid transport.
func newBareRepoWithFiles(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	bareDir := filepath.Join(root, "bare.git")
	if _, err := gogit.PlainInit(bareDir, true); err != nil {
		t.Fatalf("failed to init bare repo: %v", err)
	}
	if len(files) == 0 {
		return bareDir
	}

	workDir := filepath.Join(root, "seed-work")
	repo, err := gogit.PlainInitWithOptions(workDir, &gogit.PlainInitOptions{
		InitOptions: gogit.InitOptions{DefaultBranch: plumbing.NewBranchReferenceName("main")},
	})
	if err != nil {
		t.Fatalf("failed to init seed work repo: %v", err)
	}
	for path, content := range files {
		full := filepath.Join(workDir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatalf("failed to mkdir for seed file: %v", err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write seed file: %v", err)
		}
	}
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("failed to get seed worktree: %v", err)
	}
	if err := wt.AddGlob("."); err != nil {
		t.Fatalf("failed to add seed files: %v", err)
	}
	if _, err := wt.Commit("seed", &gogit.CommitOptions{
		Author: &object.Signature{Name: "seed", Email: "seed@test.local", When: time.Now()},
	}); err != nil {
		t.Fatalf("failed to commit seed files: %v", err)
	}
	if _, err := repo.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{bareDir}}); err != nil {
		t.Fatalf("failed to add seed remote: %v", err)
	}
	if err := repo.Push(&gogit.PushOptions{RefSpecs: []config.RefSpec{"refs/heads/main:refs/heads/main"}}); err != nil {
		t.Fatalf("failed to push seed commit: %v", err)
	}
	return bareDir
}

// newEmptyBareRepo creates a local bare git repo with no commits, standing in
// for a freshly created (empty) forge repository awaiting its first push.
func newEmptyBareRepo(t *testing.T) string {
	t.Helper()
	bareDir := filepath.Join(t.TempDir(), "bare.git")
	if _, err := gogit.PlainInit(bareDir, true); err != nil {
		t.Fatalf("failed to init bare repo: %v", err)
	}
	return bareDir
}

// readFileFromBareRepo clones bareDir's main branch into a scratch dir and
// returns the content of path within it, so tests can assert on what was
// actually pushed rather than trust the service's own success return value.
func readFileFromBareRepo(t *testing.T, bareDir, path string) (string, bool) {
	t.Helper()
	cloneDir := filepath.Join(t.TempDir(), "verify-clone")
	_, err := gogit.PlainClone(cloneDir, false, &gogit.CloneOptions{
		URL:           bareDir,
		ReferenceName: plumbing.NewBranchReferenceName("main"),
	})
	if err != nil {
		t.Fatalf("failed to clone %s for verification: %v", bareDir, err)
	}
	data, err := os.ReadFile(filepath.Join(cloneDir, path))
	if os.IsNotExist(err) {
		return "", false
	}
	if err != nil {
		t.Fatalf("failed to read %s from verification clone: %v", path, err)
	}
	return string(data), true
}

func TestAppService_CreateAppFromTemplate_E2E(t *testing.T) {
	templateRepo := newBareRepoWithFiles(t, map[string]string{"README.md": "# Template\n"})
	destRepo := newEmptyBareRepo(t)

	appRepo := mocks.NewMockAppRepository(t)
	projectRepo := mocks.NewMockProjectRepository(t)
	templateRepository := mocks.NewMockTemplateRepository(t)
	forge := forgemocks.NewMockForge(t)
	ci := cimocks.NewMockCI(t)

	templateRepository.EXPECT().GetByID(context.Background(), uint(1)).
		Return(&model.Template{ID: 1, Name: "go-template", CloneURL: templateRepo}, nil)
	appRepo.EXPECT().GetByName(context.Background(), "my-app").Return(nil, nil)
	forge.EXPECT().CreateRepo(context.Background(), "owner", "my-app").
		Return(&model.Repo{ID: 42, Name: "my-app", FullName: "owner/my-app", URL: "https://forge.example/owner/my-app", CloneURL: destRepo}, nil)
	ci.EXPECT().ActivateRepo(context.Background(), int64(42), "owner/my-app").
		Return(&model.CIRepo{ID: 99, Slug: "owner/my-app", URL: "https://ci.example/owner/my-app"}, nil)

	app := &model.App{Name: "my-app", TemplateID: ptr(uint(1))}
	appRepo.EXPECT().Save(context.Background(), app).Return(nil)

	svc := NewAppService(appRepo, projectRepo, templateRepository, forge, ci, git.New("", "test", "test@test.local", "main"), "owner", "", "secret", "main")

	got, err := svc.CreateAppFromTemplate(context.Background(), app)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.RepoID != 42 || got.RepoFullName != "owner/my-app" || got.CIID != 99 {
		t.Fatalf("unexpected app fields: %+v", got)
	}

	content, ok := readFileFromBareRepo(t, destRepo, ".abyss.yml")
	if !ok {
		t.Fatal("expected .abyss.yml to have been pushed to the new repo")
	}
	if !strings.Contains(content, "name: my-app") {
		t.Fatalf("unexpected .abyss.yml content: %s", content)
	}
	if _, ok := readFileFromBareRepo(t, destRepo, "README.md"); !ok {
		t.Fatal("expected template files to have carried over to the new repo")
	}
}

func TestAppService_CreateAppFromRepo_E2E(t *testing.T) {
	t.Run("creates .abyss.yml when missing from the repo", func(t *testing.T) {
		existingRepo := newBareRepoWithFiles(t, map[string]string{"main.go": "package main\n"})

		appRepo := mocks.NewMockAppRepository(t)
		projectRepo := mocks.NewMockProjectRepository(t)
		templateRepository := mocks.NewMockTemplateRepository(t)
		forge := forgemocks.NewMockForge(t)
		ci := cimocks.NewMockCI(t)

		appRepo.EXPECT().GetByName(context.Background(), "existing-app").Return(nil, nil)
		forge.EXPECT().GetRepo(context.Background(), int64(7)).
			Return(&model.Repo{ID: 7, Name: "existing-app", FullName: "owner/existing-app", URL: "https://forge.example/owner/existing-app", CloneURL: existingRepo}, nil)
		ci.EXPECT().ActivateRepo(context.Background(), int64(7), "owner/existing-app").
			Return(&model.CIRepo{ID: 55, Slug: "owner/existing-app", URL: "https://ci.example/owner/existing-app"}, nil)

		app := &model.App{Name: "existing-app", RepoID: 7}
		appRepo.EXPECT().Save(context.Background(), app).Return(nil)

		svc := NewAppService(appRepo, projectRepo, templateRepository, forge, ci, git.New("", "test", "test@test.local", "main"), "owner", "", "secret", "main")

		got, err := svc.CreateAppFromRepo(context.Background(), app)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.CIID != 55 {
			t.Fatalf("unexpected app: %+v", got)
		}

		content, ok := readFileFromBareRepo(t, existingRepo, ".abyss.yml")
		if !ok {
			t.Fatal("expected .abyss.yml to have been created and pushed")
		}
		if !strings.Contains(content, "name: existing-app") {
			t.Fatalf("unexpected .abyss.yml content: %s", content)
		}
	})

	t.Run("leaves an existing .abyss.yml untouched", func(t *testing.T) {
		existingRepo := newBareRepoWithFiles(t, map[string]string{
			"main.go":    "package main\n",
			".abyss.yml": "name: hand-written\ndescription: \"custom\"\n",
		})

		appRepo := mocks.NewMockAppRepository(t)
		projectRepo := mocks.NewMockProjectRepository(t)
		templateRepository := mocks.NewMockTemplateRepository(t)
		forge := forgemocks.NewMockForge(t)
		ci := cimocks.NewMockCI(t)

		appRepo.EXPECT().GetByName(context.Background(), "existing-app").Return(nil, nil)
		forge.EXPECT().GetRepo(context.Background(), int64(7)).
			Return(&model.Repo{ID: 7, Name: "existing-app", FullName: "owner/existing-app", URL: "https://forge.example/owner/existing-app", CloneURL: existingRepo}, nil)
		ci.EXPECT().ActivateRepo(context.Background(), int64(7), "owner/existing-app").
			Return(&model.CIRepo{ID: 55, Slug: "owner/existing-app", URL: "https://ci.example/owner/existing-app"}, nil)

		app := &model.App{Name: "existing-app", RepoID: 7}
		appRepo.EXPECT().Save(context.Background(), app).Return(nil)

		svc := NewAppService(appRepo, projectRepo, templateRepository, forge, ci, git.New("", "test", "test@test.local", "main"), "owner", "", "secret", "main")

		if _, err := svc.CreateAppFromRepo(context.Background(), app); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, ok := readFileFromBareRepo(t, existingRepo, ".abyss.yml")
		if !ok {
			t.Fatal("expected the existing .abyss.yml to still be present")
		}
		if content != "name: hand-written\ndescription: \"custom\"\n" {
			t.Fatalf("expected the hand-written .abyss.yml to be left untouched, got: %s", content)
		}
	})
}
