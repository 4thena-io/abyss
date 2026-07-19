package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"io"
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
	forgeType     string
}

func NewHookHandler(docsService *service.DocsService, branch, webhookSecret, forgeType string) *HookHandler {
	return &HookHandler{
		docsService:   *docsService,
		branch:        branch,
		webhookSecret: webhookSecret,
		forgeType:     forgeType,
	}
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

// verifySignature checks the webhook request's authenticity against the
// configured secret. GitHub/Gitea/Forgejo sign the raw body with
// HMAC-SHA256; GitLab instead sends the secret verbatim in a header. Both
// forms are compared in constant time to avoid timing side-channels.
func (h *HookHandler) verifySignature(r *http.Request, body []byte) bool {
	if h.webhookSecret == "" {
		return false
	}

	switch h.forgeType {
	case "gitlab":
		return subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Gitlab-Token")), []byte(h.webhookSecret)) == 1
	case "github":
		return verifyHMACSignature(r.Header.Get("X-Hub-Signature-256"), "sha256=", h.webhookSecret, body)
	default: // "gitea", "forgejo"
		sig := r.Header.Get("X-Forgejo-Signature")
		if sig == "" {
			sig = r.Header.Get("X-Gitea-Signature")
		}
		return verifyHMACSignature(sig, "", h.webhookSecret, body)
	}
}

func verifyHMACSignature(header, prefix, secret string, body []byte) bool {
	if header == "" || !strings.HasPrefix(header, prefix) {
		return false
	}
	sig, err := hex.DecodeString(strings.TrimPrefix(header, prefix))
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := mac.Sum(nil)

	return hmac.Equal(sig, expected)
}

// Forge handles push-event webhooks forwarded by the repository provider.
func (h *HookHandler) Forge(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if !h.verifySignature(r, body) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "appID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid app id", http.StatusBadRequest)
		return
	}

	var payload pushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

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
