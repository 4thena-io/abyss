package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/4thena-io/abyss/internal/git"
	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service/mocks"
)

func TestDocsService_RenderDocs(t *testing.T) {
	t.Run("errors when git client is not configured", func(t *testing.T) {
		svc := NewDocsService(mocks.NewMockAppRepository(t), nil, t.TempDir())

		err := svc.RenderDocs(context.Background(), 1)
		if err == nil {
			t.Fatal("expected an error when git client is nil")
		}
	})

	t.Run("errors when app does not exist", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(nil, nil)

		// A non-nil git client is required to reach the app lookup, but
		// RenderDocs never calls it before returning here since the app
		// isn't found.
		svc := NewDocsService(appRepo, git.New("", "", "", ""), t.TempDir())

		err := svc.RenderDocs(context.Background(), 1)
		if err == nil {
			t.Fatal("expected an error when app is not found")
		}
	})
}

func TestDocsService_GetDocs(t *testing.T) {
	t.Run("returns ErrNotFound when app does not exist", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(nil, nil)

		svc := NewDocsService(appRepo, nil, t.TempDir())

		_, err := svc.GetDocs(context.Background(), 1)
		if err != ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrNotFound when docs have not been rendered yet", func(t *testing.T) {
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(&model.App{ID: 1}, nil)

		svc := NewDocsService(appRepo, nil, t.TempDir())

		_, err := svc.GetDocs(context.Background(), 1)
		if err != ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns previously stored docs", func(t *testing.T) {
		docsDir := t.TempDir()
		appRepo := mocks.NewMockAppRepository(t)
		appRepo.EXPECT().GetByID(context.Background(), uint(1)).Return(&model.App{ID: 1}, nil)

		svc := NewDocsService(appRepo, nil, docsDir)

		rendered := &RenderedDocs{
			Pages: []DocPage{{Path: "index.md", HTML: "<h1>Hi</h1>"}},
			Nav:   []NavNode{{Title: "Overview", Path: "index.md"}},
		}
		if err := svc.store(1, rendered); err != nil {
			t.Fatalf("failed to seed stored docs: %v", err)
		}

		got, err := svc.GetDocs(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Pages) != 1 || got.Pages[0].HTML != "<h1>Hi</h1>" {
			t.Fatalf("unexpected docs: %+v", got)
		}
	})
}

func TestDocsService_renderDocsDir(t *testing.T) {
	svc := NewDocsService(nil, nil, t.TempDir())

	t.Run("returns empty result when docs folder is missing", func(t *testing.T) {
		got, err := svc.renderDocsDir(filepath.Join(t.TempDir(), "does-not-exist"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Pages) != 0 {
			t.Fatalf("expected no pages, got %+v", got)
		}
	})

	t.Run("renders markdown files, ignoring non-.md files", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "index.md"), []byte("# Getting Started\n"), 0644); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, "guide.md"), []byte("# Guide\n"), 0644); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, "notes.txt"), []byte("not markdown"), 0644); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}

		got, err := svc.renderDocsDir(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Pages) != 2 {
			t.Fatalf("expected 2 markdown pages (non-.md file excluded), got %d: %+v", len(got.Pages), got.Pages)
		}
	})

	t.Run("root index.md becomes a leading Overview nav node", func(t *testing.T) {
		root := t.TempDir()
		mustWrite(t, filepath.Join(root, "index.md"), "# Home\n")
		mustWrite(t, filepath.Join(root, "guide.md"), "# Guide\n")

		got, err := svc.renderDocsDir(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Nav) != 2 {
			t.Fatalf("expected 2 nav nodes, got %d: %+v", len(got.Nav), got.Nav)
		}
		if got.Nav[0].Title != "Overview" || got.Nav[0].Path != "index.md" {
			t.Fatalf("expected root index.md to be a leading Overview node, got %+v", got.Nav[0])
		}
	})

	t.Run("nested directories preserve full depth", func(t *testing.T) {
		root := t.TempDir()
		mustWrite(t, filepath.Join(root, "a", "b", "c.md"), "# Deep\n")

		got, err := svc.renderDocsDir(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Nav) != 1 || got.Nav[0].Title != "A" || len(got.Nav[0].Children) != 1 {
			t.Fatalf("expected top-level 'A' section, got %+v", got.Nav)
		}
		b := got.Nav[0].Children[0]
		if b.Title != "B" || len(b.Children) != 1 {
			t.Fatalf("expected nested 'B' section under 'A', got %+v", b)
		}
		c := b.Children[0]
		if c.Title != "C" || c.Path != "a/b/c.md" {
			t.Fatalf("expected leaf 'C' at full depth path, got %+v", c)
		}
	})

	t.Run("a directory's index.md becomes its section header, not a separate Overview row", func(t *testing.T) {
		root := t.TempDir()
		mustWrite(t, filepath.Join(root, "getting-started", "index.md"), "# Getting Started\n")
		mustWrite(t, filepath.Join(root, "getting-started", "installation.md"), "# Installation\n")

		got, err := svc.renderDocsDir(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Nav) != 1 {
			t.Fatalf("expected 1 top-level section, got %d: %+v", len(got.Nav), got.Nav)
		}
		section := got.Nav[0]
		if section.Title != "Getting Started" {
			t.Fatalf("expected section title 'Getting Started', got %q", section.Title)
		}
		if section.Path != "getting-started/index.md" {
			t.Fatalf("expected section header to link to its index.md, got path %q", section.Path)
		}
		if len(section.Children) != 1 || section.Children[0].Path != "getting-started/installation.md" {
			t.Fatalf("expected exactly one sibling child (installation.md), no duplicate Overview row, got %+v", section.Children)
		}
	})

	t.Run("a present nav.yml overrides auto-discovered titles and order", func(t *testing.T) {
		root := t.TempDir()
		mustWrite(t, filepath.Join(root, "api", "index.md"), "# API\n")
		mustWrite(t, filepath.Join(root, "getting-started", "index.md"), "# Getting Started\n")
		mustWrite(t, filepath.Join(root, "nav.yml"), `
- title: Getting Started
  path: getting-started
- title: API Reference
  path: api
`)

		got, err := svc.renderDocsDir(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Nav) != 2 || got.Nav[0].Title != "Getting Started" || got.Nav[1].Title != "API Reference" {
			t.Fatalf("expected nav.yml order/titles to apply, got %+v", got.Nav)
		}
		// nav.yml itself must not surface as a page.
		for _, p := range got.Pages {
			if p.Path == "nav.yml" {
				t.Fatalf("expected nav.yml to be excluded from Pages, got %+v", got.Pages)
			}
		}
	})

	t.Run("the root Overview stays first regardless of nav.yml, since it can't be addressed by a manifest entry", func(t *testing.T) {
		root := t.TempDir()
		mustWrite(t, filepath.Join(root, "index.md"), "# Home\n")
		mustWrite(t, filepath.Join(root, "api", "index.md"), "# API\n")
		mustWrite(t, filepath.Join(root, "getting-started", "index.md"), "# Getting Started\n")
		mustWrite(t, filepath.Join(root, "nav.yml"), `
- title: Getting Started
  path: getting-started
- title: API Reference
  path: api
`)

		got, err := svc.renderDocsDir(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Nav) != 3 || got.Nav[0].Title != "Overview" || got.Nav[0].Path != "index.md" {
			t.Fatalf("expected Overview to stay pinned first, got %+v", got.Nav)
		}
		if got.Nav[1].Title != "Getting Started" || got.Nav[2].Title != "API Reference" {
			t.Fatalf("expected manifest order to apply to the remaining entries, got %+v", got.Nav)
		}
	})

	t.Run("a malformed nav.yml is non-fatal: rendering still succeeds with the auto-discovered nav", func(t *testing.T) {
		root := t.TempDir()
		mustWrite(t, filepath.Join(root, "getting-started", "index.md"), "# Getting Started\n")
		mustWrite(t, filepath.Join(root, "nav.yml"), "title: [not a list of entries\n")

		got, err := svc.renderDocsDir(root)
		if err != nil {
			t.Fatalf("expected malformed nav.yml to be non-fatal, got error: %v", err)
		}
		if len(got.Nav) != 1 || got.Nav[0].Title != "Getting Started" {
			t.Fatalf("expected auto-discovered nav to be used as a fallback, got %+v", got.Nav)
		}
	})
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create fixture dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}
}

func TestDocsService_renderMarkdown_EscapesRawHTML(t *testing.T) {
	svc := NewDocsService(nil, nil, t.TempDir())
	md := svc.newMarkdown()

	html, err := svc.renderMarkdown(md, []byte("hello <script>alert(1)</script> world\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Without WithUnsafe(), goldmark doesn't pass raw HTML tags through in
	// any form (escaped or otherwise) — it drops them and leaves an HTML
	// comment in their place, so <script> can never reach the browser.
	if strings.Contains(html, "<script>") || strings.Contains(html, "</script>") {
		t.Fatalf("expected raw <script> tags not to appear in the output at all, got: %s", html)
	}
	if !strings.Contains(html, "raw HTML omitted") {
		t.Fatalf("expected goldmark's raw-HTML placeholder in output, got: %s", html)
	}
}
