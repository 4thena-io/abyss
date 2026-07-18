package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/4thena-io/abyss/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type HookHandler struct {
	docsService   service.DocsService
	branch        string
	webhookSecret string
}

func NewHookHandler(docsService *service.DocsService, branch, webhookSecret string) *HookHandler {
	return &HookHandler{*docsService, branch, webhookSecret}
}

// pushPayload captures the common push event shape across GitHub, GitLab, Gitea, and Forgejo.
type pushPayload struct {
	Ref     string `json:"ref"` // e.g. "refs/heads/main"
	Commits []struct {
		Added    []string `json:"added"`
		Modified []string `json:"modified"`
		Removed  []string `json:"removed"`
	} `json:"commits"`
}

func (p pushPayload) isOnBranch(branch string) bool {
	return p.Ref == "refs/heads/"+branch
}

func (p pushPayload) touchesDocs() bool {
	// If the forge sends no file lists at all (empty commits or missing diff info),
	// we can't tell what changed — assume docs might be affected.
	hasAnyFiles := false
	for _, c := range p.Commits {
		if len(c.Added)+len(c.Modified)+len(c.Removed) > 0 {
			hasAnyFiles = true
			break
		}
	}
	if !hasAnyFiles {
		return true
	}

	for _, c := range p.Commits {
		for _, f := range append(append(c.Added, c.Modified...), c.Removed...) {
			if strings.HasPrefix(f, "docs/") || f == ".abyss.yml" {
				return true
			}
		}
	}
	return false
}

// Forge handles push-event webhooks forwarded by the repository provider.
func (h *HookHandler) Forge(w http.ResponseWriter, r *http.Request) {
	if h.webhookSecret != "" && r.URL.Query().Get("access_token") != h.webhookSecret {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "appID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid app id", http.StatusBadRequest)
		return
	}

	var payload pushPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()

	logger := log.With().Uint64("app_id", id).Str("ref", payload.Ref).Logger()
	logger.Debug().Msg("forge webhook received")

	if !payload.isOnBranch(h.branch) {
		logger.Info().Str("want_branch", h.branch).Msg("webhook skipped: wrong branch")
		w.WriteHeader(http.StatusOK)
		return
	}

	if !payload.touchesDocs() {
		logger.Info().Msg("webhook skipped: no docs changes")
		w.WriteHeader(http.StatusOK)
		return
	}

	logger.Info().Msg("docs change detected, rendering")

	if err := h.docsService.RenderDocs(r.Context(), uint(id)); err != nil {
		logger.Error().Err(err).Msg("failed to render docs")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	logger.Info().Msg("docs rendered successfully")
	w.WriteHeader(http.StatusOK)
}
