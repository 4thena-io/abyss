package service

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// navManifestEntry mirrors one entry of an optional docs/nav.yml manifest,
// letting doc authors override the auto-discovered nav's title and ordering
// without renaming files.
//
// path is matched against a NavNode's navNodeKey: a directory's raw (slug)
// name for sections (e.g. "getting-started"), or the file's basename
// including extension for leaf pages (e.g. "installation.md") — the same
// two forms already visible on disk, so authors don't need to think about
// how titles get humanized.
//
// A leaf item may be written as a bare string ("installation.md") when no
// title override is needed, or as a mapping when it is.
type navManifestEntry struct {
	Title string
	Path  string
	Items []navManifestEntry
}

// UnmarshalYAML accepts either a bare scalar (just a path, e.g.
// "installation.md") or a full mapping with title/path/items.
func (e *navManifestEntry) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		e.Path = value.Value
		return nil
	}

	var mapped struct {
		Title string             `yaml:"title"`
		Path  string             `yaml:"path"`
		Items []navManifestEntry `yaml:"items"`
	}
	if err := value.Decode(&mapped); err != nil {
		return err
	}
	e.Title, e.Path, e.Items = mapped.Title, mapped.Path, mapped.Items
	return nil
}

// loadNavManifest reads docs/nav.yml under docsRoot, if present. A missing
// file returns (nil, nil); a malformed file returns a non-nil error that
// callers must treat as non-fatal (log and fall back to the auto-discovered
// nav), never as a reason to fail rendering.
func loadNavManifest(docsRoot string) ([]navManifestEntry, error) {
	data, err := os.ReadFile(filepath.Join(docsRoot, "nav.yml"))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var entries []navManifestEntry
	if err := yaml.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// applyNavManifest overlays manifest ordering/titles onto an auto-discovered
// nav tree. It never removes or hides anything: a manifest entry matched
// against an auto-discovered node (see navNodeKey) takes the manifest's
// title and is emitted in manifest order; any auto-discovered sibling not
// mentioned in the manifest is appended afterward, in its existing
// (alphabetical) order. A manifest entry with no matching on-disk node is
// silently skipped — a stale manifest entry is a no-op, not an error.
func applyNavManifest(auto []NavNode, manifest []navManifestEntry) []NavNode {
	byKey := make(map[string]NavNode, len(auto))
	for _, node := range auto {
		byKey[navNodeKey(node)] = node
	}

	used := make(map[string]bool, len(auto))
	ordered := make([]NavNode, 0, len(auto))
	for _, entry := range manifest {
		node, ok := byKey[entry.Path]
		if !ok {
			continue
		}
		used[entry.Path] = true
		if entry.Title != "" {
			node.Title = entry.Title
		}
		if len(entry.Items) > 0 && len(node.Children) > 0 {
			node.Children = applyNavManifest(node.Children, entry.Items)
		}
		ordered = append(ordered, node)
	}

	for _, node := range auto {
		if !used[navNodeKey(node)] {
			ordered = append(ordered, node)
		}
	}

	return ordered
}

// navNodeKey is the on-disk name a manifest entry's path matches against: a
// section's raw directory name (its slug, since Title is already humanized
// and Path, if set, points at the section's index.md rather than the
// directory), or a leaf page's filename including extension.
func navNodeKey(n NavNode) string {
	if n.slug != "" {
		return n.slug
	}
	return filepath.Base(n.Path)
}
