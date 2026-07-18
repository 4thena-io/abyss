package service

import (
	"context"
	"os"
	"path/filepath"
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
			Pages:     []DocPage{{Path: "index.md", HTML: "<h1>Hi</h1>"}},
			Structure: []string{"index.md"},
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

	t.Run("renders markdown files and extracts TOC from index.md", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "index.md"), []byte("# Getting Started\n\ntext\n\n## Install\n\nmore text\n"), 0644); err != nil {
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
		if len(got.TOC) != 2 {
			t.Fatalf("expected 2 TOC entries from index.md headings, got %d: %+v", len(got.TOC), got.TOC)
		}
		if got.TOC[0].Title != "Getting Started" || got.TOC[0].Level != 1 {
			t.Fatalf("unexpected first TOC entry: %+v", got.TOC[0])
		}
		if got.TOC[1].Title != "Install" || got.TOC[1].Level != 2 {
			t.Fatalf("unexpected second TOC entry: %+v", got.TOC[1])
		}
		if got.TOC[0].Anchor == "" || got.TOC[1].Anchor == "" {
			t.Fatalf("expected auto-generated heading anchors, got %+v", got.TOC)
		}
	})
}

func TestDocsService_extractTOC_nestedInlineFormatting(t *testing.T) {
	svc := NewDocsService(nil, nil, t.TempDir())
	md := svc.newMarkdown()
	source := []byte("# Hello **bold** world\n")

	toc := svc.extractTOC(md, source)
	if len(toc) != 1 || toc[0].Title != "Hello bold world" {
		t.Fatalf("expected heading text to flatten inline formatting, got %+v", toc)
	}
}
