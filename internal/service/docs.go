package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/4thena-io/abyss/internal/git"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type DocPage struct {
	Path string `json:"path"`
	HTML string `json:"html"`
}

// NavNode is one entry in the docs sidebar tree. Children is deliberately
// NOT omitempty: a directory node always serializes it (as `[]` when the
// directory has no other files, `null` for a plain leaf page) so the
// frontend can tell "this is a folder" apart from "this is a leaf" without
// that depending on whether the folder happens to have sibling files. If it
// also has Path set, the section header links to that directory's index.md
// (folder-as-clickable-header); without Path it's a non-clickable group
// label (no index.md in that directory).
type NavNode struct {
	Title    string    `json:"title"`
	Path     string    `json:"path,omitempty"`
	Children []NavNode `json:"children"`

	// slug is the raw (non-humanized) directory name for section nodes; it
	// is unexported (never reaches the JSON contract) and exists purely so
	// applyNavManifest can match a nav.yml entry's `path: getting-started`
	// against this node even though Title has already been humanized to
	// "Getting Started" and Path (when set) points at the section's
	// index.md rather than the directory itself.
	slug string
}

type RenderedDocs struct {
	Pages []DocPage `json:"pages"`
	Nav   []NavNode `json:"nav"`
}

type DocsService struct {
	appRepository AppRepository
	git           *git.GitClient
	docsDir       string
}

func NewDocsService(appRepository AppRepository, gitClient *git.GitClient, docsDir string) *DocsService {
	return &DocsService{
		appRepository: appRepository,
		git:           gitClient,
		docsDir:       docsDir,
	}
}

func (s *DocsService) RenderDocs(ctx context.Context, appID uint) error {
	logger := log.With().Uint("app_id", appID).Logger()

	if s.git == nil {
		return fmt.Errorf("git client not configured")
	}

	app, err := s.appRepository.GetByID(ctx, appID)
	if err != nil {
		return fmt.Errorf("failed to get app: %w", err)
	}
	if app == nil {
		return fmt.Errorf("app not found")
	}

	logger = logger.With().Str("repo", app.RepoFullName).Logger()

	tmpDir, err := os.MkdirTemp("", "abyss-docs-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	logger.Debug().Str("tmp_dir", tmpDir).Msg("cloning docs")
	if err := s.git.CloneSubset(app.CloneURL, tmpDir, []string{"docs"}); err != nil {
		return fmt.Errorf("failed to clone docs: %w", err)
	}

	docsPath := filepath.Join(tmpDir, "docs")

	logger.Debug().Msg("rendering markdown")
	rendered, err := s.renderDocsDir(docsPath)
	if err != nil {
		return fmt.Errorf("failed to render docs dir: %w", err)
	}

	logger.Debug().Int("pages", len(rendered.Pages)).Int("nav_nodes", len(rendered.Nav)).Msg("storing rendered docs")
	if err := s.store(appID, rendered); err != nil {
		return err
	}

	logger.Info().Int("pages", len(rendered.Pages)).Str("output", filepath.Join(s.docsDir, fmt.Sprintf("%d", appID))).Msg("docs stored")
	return nil
}

func (s *DocsService) store(appID uint, rendered *RenderedDocs) error {
	appDocsDir := filepath.Join(s.docsDir, fmt.Sprintf("%d", appID))
	if err := os.MkdirAll(appDocsDir, 0755); err != nil {
		return fmt.Errorf("failed to create docs dir: %w", err)
	}

	data, err := json.Marshal(rendered)
	if err != nil {
		return fmt.Errorf("failed to marshal rendered docs: %w", err)
	}

	return os.WriteFile(filepath.Join(appDocsDir, "docs.json"), data, 0644)
}

func (s *DocsService) newMarkdown() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		// No WithUnsafe(): raw HTML in repo-authored markdown (e.g. a
		// malicious docs/*.md from anyone with push access) is escaped
		// rather than passed through, since the rendered output is later
		// injected into the browser via v-html.
	)
}

func (s *DocsService) GetDocs(ctx context.Context, appID uint) (*RenderedDocs, error) {
	app, err := s.appRepository.GetByID(ctx, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get app: %w", err)
	}
	if app == nil {
		return nil, ErrNotFound
	}

	data, err := os.ReadFile(filepath.Join(s.docsDir, fmt.Sprintf("%d", appID), "docs.json"))
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read docs: %w", err)
	}

	var docs RenderedDocs
	if err := json.Unmarshal(data, &docs); err != nil {
		return nil, fmt.Errorf("failed to parse docs: %w", err)
	}

	return &docs, nil
}

func (s *DocsService) renderDocsDir(root string) (*RenderedDocs, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Warn().Str("path", root).Msg("docs folder not found in repo, skipping render")
		return &RenderedDocs{}, nil
	}

	md := s.newMarkdown()

	var pages []DocPage

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		rendered, err := s.renderMarkdown(md, content)
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)

		pages = append(pages, DocPage{Path: rel, HTML: rendered})

		return nil
	})
	if err != nil {
		return nil, err
	}

	nav := buildNavTree(pages)
	if manifest, err := loadNavManifest(root); err != nil {
		log.Warn().Err(err).Str("path", root).Msg("failed to parse docs/nav.yml, falling back to auto-discovered nav")
	} else if manifest != nil {
		// The root's own Overview (docs/index.md) is pinned first by
		// buildNavTree and isn't something a manifest entry can address
		// (there's no directory name or filename to write as its `path`) —
		// exclude it from the overlay pass so it can't get shuffled to the
		// end as just another unmatched node.
		if len(nav) > 0 && nav[0].Path == "index.md" && nav[0].Children == nil {
			nav = append([]NavNode{nav[0]}, applyNavManifest(nav[1:], manifest)...)
		} else {
			nav = applyNavManifest(nav, manifest)
		}
	}

	return &RenderedDocs{Pages: pages, Nav: nav}, nil
}

func (s *DocsService) renderMarkdown(md goldmark.Markdown, content []byte) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// navDir accumulates one directory's leaf pages and subdirectories while
// buildNavTree groups the flat page list into a tree.
type navDir struct {
	title    string
	index    string // rel path of this dir's index.md, "" if none
	children []NavNode
	dirs     map[string]*navDir
}

// buildNavTree groups a flat, walk-ordered page list into a directory tree.
// A directory's index.md becomes its own NavNode.Path (the folder-as-header
// case) instead of also appearing as a sibling leaf page.
func buildNavTree(pages []DocPage) []NavNode {
	root := &navDir{dirs: map[string]*navDir{}}

	getOrCreateDir := func(parent *navDir, name string) *navDir {
		if d, ok := parent.dirs[name]; ok {
			return d
		}
		d := &navDir{title: humanizeTitle(name), dirs: map[string]*navDir{}}
		parent.dirs[name] = d
		return d
	}

	for _, p := range pages {
		parts := strings.Split(p.Path, "/")
		dir := root
		for _, seg := range parts[:len(parts)-1] {
			dir = getOrCreateDir(dir, seg)
		}

		file := parts[len(parts)-1]
		if file == "index.md" {
			dir.index = p.Path
			continue
		}
		dir.children = append(dir.children, NavNode{
			Title: humanizeTitle(strings.TrimSuffix(file, ".md")),
			Path:  p.Path,
		})
	}

	nodes := navDirNodes(root)
	if root.index != "" {
		nodes = append([]NavNode{{Title: "Overview", Path: root.index}}, nodes...)
	}
	return nodes
}

func navDirNodes(d *navDir) []NavNode {
	nodes := append([]NavNode{}, d.children...)

	names := make([]string, 0, len(d.dirs))
	for name := range d.dirs {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		sub := d.dirs[name]
		node := NavNode{Title: sub.title, Children: navDirNodes(sub), slug: name}
		if sub.index != "" {
			node.Path = sub.index
		}
		nodes = append(nodes, node)
	}

	sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].Title < nodes[j].Title })
	return nodes
}

// humanizeTitle turns a file or directory name into display text:
// "getting-started" -> "Getting Started".
func humanizeTitle(name string) string {
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	words := strings.Fields(name)
	for i, w := range words {
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	return strings.Join(words, " ")
}
