package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/4thena-io/abyss/internal/git"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

type DocPage struct {
	Path string `json:"path"`
	HTML string `json:"html"`
}

type TOCItem struct {
	Title  string `json:"title"`
	Anchor string `json:"anchor"`
	Level  int    `json:"level"`
}

type RenderedDocs struct {
	Pages     []DocPage `json:"pages"`
	TOC       []TOCItem `json:"toc"`
	Structure []string  `json:"structure"`
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

	logger.Debug().Int("pages", len(rendered.Pages)).Int("toc_items", len(rendered.TOC)).Msg("storing rendered docs")
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
		goldmark.WithRendererOptions(
			goldmarkhtml.WithUnsafe(),
		),
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
	var structure []string
	var toc []TOCItem

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

		pages = append(pages, DocPage{Path: rel, HTML: rendered})
		structure = append(structure, rel)

		if filepath.Base(path) == "index.md" {
			toc = s.extractTOC(md, content)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &RenderedDocs{Pages: pages, TOC: toc, Structure: structure}, nil
}

func (s *DocsService) renderMarkdown(md goldmark.Markdown, content []byte) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *DocsService) extractTOC(md goldmark.Markdown, source []byte) []TOCItem {
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	var toc []TOCItem

	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		h, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		anchor := ""
		if id, ok := h.AttributeString("id"); ok {
			anchor = string(id.([]byte))
		}

		toc = append(toc, TOCItem{
			Title:  string(headingText(h, source)),
			Anchor: anchor,
			Level:  h.Level,
		})

		return ast.WalkContinue, nil
	})

	return toc
}

// headingText concatenates the plain text of a heading's descendants.
// ast.Node.Text is deprecated; this walks *ast.Text leaves directly instead.
func headingText(n ast.Node, source []byte) []byte {
	var buf bytes.Buffer
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			buf.Write(t.Segment.Value(source))
			continue
		}
		buf.Write(headingText(c, source))
	}
	return buf.Bytes()
}
