package service

import (
	"path/filepath"
	"testing"
)

func TestLoadNavManifest(t *testing.T) {
	t.Run("returns nil, nil when nav.yml is absent", func(t *testing.T) {
		entries, err := loadNavManifest(t.TempDir())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if entries != nil {
			t.Fatalf("expected nil entries, got %+v", entries)
		}
	})

	t.Run("parses valid YAML, including bare-string leaf items", func(t *testing.T) {
		root := t.TempDir()
		mustWrite(t, filepath.Join(root, "nav.yml"), `
- title: Getting Started
  path: getting-started
  items:
    - installation.md
    - configuration.md
- title: API Reference
  path: api
`)

		entries, err := loadNavManifest(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(entries) != 2 {
			t.Fatalf("expected 2 top-level entries, got %d: %+v", len(entries), entries)
		}
		if entries[0].Title != "Getting Started" || entries[0].Path != "getting-started" {
			t.Fatalf("unexpected first entry: %+v", entries[0])
		}
		if len(entries[0].Items) != 2 || entries[0].Items[0].Path != "installation.md" {
			t.Fatalf("unexpected items: %+v", entries[0].Items)
		}
	})

	t.Run("returns an error for malformed YAML", func(t *testing.T) {
		root := t.TempDir()
		mustWrite(t, filepath.Join(root, "nav.yml"), "title: [this is not a list of entries\n")

		if _, err := loadNavManifest(root); err == nil {
			t.Fatal("expected an error for malformed YAML")
		}
	})
}

func TestApplyNavManifest(t *testing.T) {
	auto := []NavNode{
		{Title: "Api", slug: "api"},
		{Title: "Getting Started", Path: "getting-started/index.md", slug: "getting-started", Children: []NavNode{
			{Title: "Configuration", Path: "getting-started/configuration.md"},
			{Title: "Installation", Path: "getting-started/installation.md"},
		}},
	}

	t.Run("reorders and retitles matched entries, preserving unmatched siblings", func(t *testing.T) {
		manifest := []navManifestEntry{
			{Title: "Getting Started", Path: "getting-started", Items: []navManifestEntry{
				{Path: "installation.md"},
				{Path: "configuration.md"},
			}},
			{Title: "API Reference", Path: "api"},
		}

		got := applyNavManifest(auto, manifest)
		if len(got) != 2 {
			t.Fatalf("expected 2 nodes, got %d: %+v", len(got), got)
		}
		if got[0].Title != "Getting Started" {
			t.Fatalf("expected 'Getting Started' first, got %+v", got[0])
		}
		if got[0].Children[0].Path != "getting-started/installation.md" {
			t.Fatalf("expected installation.md first among children, got %+v", got[0].Children)
		}
		if got[1].Title != "API Reference" {
			t.Fatalf("expected 'api' node retitled to 'API Reference', got %+v", got[1])
		}
	})

	t.Run("a manifest entry pointing at a nonexistent path is a no-op", func(t *testing.T) {
		manifest := []navManifestEntry{
			{Title: "Ghost", Path: "does-not-exist"},
		}

		got := applyNavManifest(auto, manifest)
		if len(got) != len(auto) {
			t.Fatalf("expected stale entry to be skipped, got %+v", got)
		}
	})

	t.Run("on-disk items missing from the manifest are appended after explicit ones", func(t *testing.T) {
		manifest := []navManifestEntry{
			{Path: "getting-started"},
		}

		got := applyNavManifest(auto, manifest)
		if len(got) != 2 {
			t.Fatalf("expected 2 nodes (1 explicit + 1 auto-appended), got %+v", got)
		}
		if got[0].slug != "getting-started" || got[1].slug != "api" {
			t.Fatalf("expected 'api' auto-appended after explicit 'getting-started', got %+v", got)
		}
	})
}
