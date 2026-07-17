package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/auth"
	"github.com/4thena-io/abyss/internal/model"
)

type contextKey string

const ContextKeyUser contextKey = "user"

type AuthProvider interface {
	ValidateJWT(tokenStr string) (*auth.Claims, error)
	GetUserByToken(ctx context.Context, token string) (*model.User, error)
}

// ClaimsFromContext extracts the authenticated user's claims from the request context.
func ClaimsFromContext(ctx context.Context) *auth.Claims {
	claims, _ := ctx.Value(ContextKeyUser).(*auth.Claims)
	return claims
}

// RequireSetup redirects to /setup if OAuth credentials are not configured.
func RequireSetup(clientID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if clientID == "" {
				http.Redirect(w, r, "/setup", http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAuth authenticates requests via Bearer PAT or JWT session cookie.
// Injects *auth.Claims into the request context under ContextKeyUser.
// If svc is nil (instance not yet configured), all requests are rejected.
func RequireAuth(svc AuthProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if svc == nil {
				response.Unauthorized(w)
				return
			}

			var claims *auth.Claims

			// Resolve PAT from Bearer header or ?access_token query param.
			pat := ""
			if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
				pat = strings.TrimPrefix(h, "Bearer ")
			} else if q := r.URL.Query().Get("access_token"); q != "" {
				pat = q
			}

			if pat != "" {
				user, err := svc.GetUserByToken(r.Context(), pat)
				if err != nil || user == nil {
					response.Unauthorized(w)
					return
				}
				claims = &auth.Claims{
					UserID:   user.ID,
					Username: user.Username,
					IsAdmin:  user.IsAdmin,
				}
			} else {
				cookie, err := r.Cookie("session")
				if err != nil {
					response.Unauthorized(w)
					return
				}
				claims, err = svc.ValidateJWT(cookie.Value)
				if err != nil {
					response.Unauthorized(w)
					return
				}
			}

			ctx := context.WithValue(r.Context(), ContextKeyUser, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
