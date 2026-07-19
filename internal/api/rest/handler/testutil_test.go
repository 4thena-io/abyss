package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/4thena-io/abyss/internal/api/rest/middleware"
	"github.com/4thena-io/abyss/internal/auth"
	"github.com/go-chi/chi/v5"
)

// requestWithParams builds a request carrying chi URL params and, optionally,
// authenticated claims in the context, mirroring what the router/middleware
// would inject before the handler runs.
func requestWithParams(method, target string, body string, claims *auth.Claims, params map[string]string) *http.Request {
	var r *http.Request
	if body != "" {
		r = httptest.NewRequest(method, target, strings.NewReader(body))
	} else {
		r = httptest.NewRequest(method, target, nil)
	}

	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	if claims != nil {
		ctx = context.WithValue(ctx, middleware.ContextKeyUser, claims)
	}
	return r.WithContext(ctx)
}
