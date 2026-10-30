package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDatabaseAvailabilityRecoversWithoutClearingSession(t *testing.T) {
	online := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })
	h := DatabaseAvailable(func(context.Context) error {
		if !online {
			return errors.New("down")
		}
		return nil
	})(next)
	for _, tc := range []struct {
		path string
		want int
	}{{"/api/me", 503}, {"/auth/callback", 503}, {"/api/live", 204}, {"/api/ready", 204}, {"/api/config", 204}, {"/assets/app.js", 204}} {
		r := httptest.NewRecorder()
		h.ServeHTTP(r, httptest.NewRequest("GET", tc.path, nil))
		if r.Code != tc.want {
			t.Errorf("%s: %d", tc.path, r.Code)
		}
		if len(r.Result().Cookies()) != 0 {
			t.Fatal("outage must not clear session/state")
		}
	}
	online = true
	r := httptest.NewRecorder()
	h.ServeHTTP(r, httptest.NewRequest("GET", "/api/me", nil))
	if r.Code != 204 {
		t.Fatal(r.Code)
	}
}
