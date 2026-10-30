package cmd

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/4thena-io/abyss/internal/api/rest/handler"
	"github.com/4thena-io/abyss/internal/database"
	"github.com/rs/zerolog/log"
)

// startupHandler keeps liveness and static assets available while dependencies
// start. The complete router is published atomically, never partially wired.
type startupHandler struct {
	mu    sync.RWMutex
	ready http.Handler
	ui    *handler.UIHandler
}

func (h *startupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	ready := h.ready
	h.mu.RUnlock()
	if ready != nil {
		ready.ServeHTTP(w, r)
		return
	}
	if r.URL.Path == "/api/live" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/auth/") {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"message":"Database unavailable or application starting. Please retry shortly."}`))
		return
	}
	h.ui.Serve(w, r)
}

func (h *startupHandler) initialize(ctx context.Context, build func() (*Server, error)) error {
	delay := time.Second
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		srv, err := build()
		if err == nil {
			h.mu.Lock()
			h.ready = srv.http.Handler
			h.mu.Unlock()
			return nil
		}
		if !database.IsTransient(err) {
			return err
		}
		// Do not log driver errors here: some drivers include connection credentials.
		log.Warn().Dur("retry_in", delay).Msg("database unavailable; retrying initialization")
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
		if delay < 15*time.Second {
			delay *= 2
			if delay > 15*time.Second {
				delay = 15 * time.Second
			}
		}
	}
}
