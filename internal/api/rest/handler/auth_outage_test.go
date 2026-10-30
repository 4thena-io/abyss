package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/4thena-io/abyss/internal/integration/forge"
	"github.com/4thena-io/abyss/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func TestOAuthFailureDiagnosticsAreSafe(t *testing.T) {
	previous := log.Logger
	defer func() { log.Logger = previous }()
	var logs bytes.Buffer
	log.Logger = zerolog.New(&logs)
	h := &AuthHandler{}
	for _, tc := range []struct {
		stage string
		err   error
		want  int
	}{
		{"identity", forge.ErrUnavailable, 503},
		{"user_sync", &service.AuthStageError{Stage: "membership", Err: errors.New("secret-token")}, 503},
		{"user_sync", &service.AuthStageError{Stage: "database", Err: context.DeadlineExceeded}, 503},
		{"user_sync", &service.AuthStageError{Stage: "database", Err: errors.New("secret-token")}, 500},
		{"code_exchange", &oauth2.RetrieveError{Response: &http.Response{StatusCode: 400}, Body: []byte("secret-token")}, 400},
		{"code_exchange", &oauth2.RetrieveError{Response: &http.Response{StatusCode: 502}, Body: []byte("secret-token")}, 503},
	} {
		r := httptest.NewRecorder()
		h.callbackFailure(r, httptest.NewRequest("GET", "/auth/callback", nil), tc.stage, tc.err)
		if r.Code != tc.want {
			t.Errorf("%s: got %d want %d", tc.stage, r.Code, tc.want)
		}
		if strings.Contains(r.Body.String(), "secret-token") {
			t.Fatal("secret leaked in response")
		}
	}
	if strings.Contains(logs.String(), "secret-token") {
		t.Fatal("secret leaked in logs")
	}
	if !strings.Contains(logs.String(), `"stage":"membership"`) {
		t.Fatal("missing stage")
	}
	r := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/callback", nil)
	req.Header.Set("Accept", "text/html")
	h.callbackFailure(r, req, "identity", forge.ErrUnavailable)
	if r.Code != 503 || !strings.Contains(r.Body.String(), "/auth/login") {
		t.Fatal("missing retry link")
	}
}
