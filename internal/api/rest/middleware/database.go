package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/4thena-io/abyss/internal/api/rest/response"
)

// DatabaseAvailable rejects database-dependent requests without consuming OAuth
// state or clearing sessions. PingContext also lets database/sql reconnect after
// an outage. Liveness, readiness, public configuration and static UI stay usable.
func DatabaseAvailable(ping func(context.Context) error) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/live", "/api/ready", "/api/health", "/api/config", "/api/setup/status":
				next.ServeHTTP(w, r)
				return
			}
			if !strings.HasPrefix(r.URL.Path, "/api/") && !strings.HasPrefix(r.URL.Path, "/auth/") {
				next.ServeHTTP(w, r)
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), time.Second)
			err := ping(ctx)
			cancel()
			if err != nil {
				response.ServiceUnavailable(w, "Database temporarily unavailable. Please retry shortly.")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
