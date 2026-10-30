package cmd

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStartupRecovers(t *testing.T) {
	h := &startupHandler{}
	r := httptest.NewRecorder()
	h.ServeHTTP(r, httptest.NewRequest("GET", "/api/me", nil))
	if r.Code != 503 {
		t.Fatalf("got %d", r.Code)
	}
	attempts := 0
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := h.initialize(ctx, func() (*Server, error) {
		attempts++
		if attempts == 1 {
			return nil, &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("unavailable")}
		}
		return newServer(Config{}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	r = httptest.NewRecorder()
	h.ServeHTTP(r, httptest.NewRequest("GET", "/api/me", nil))
	if r.Code != 204 || attempts != 2 {
		t.Fatalf("status=%d attempts=%d", r.Code, attempts)
	}
}
func TestStartupDoesNotRetryPermanentError(t *testing.T) {
	h := &startupHandler{}
	want := errors.New("invalid configuration")
	attempts := 0
	err := h.initialize(context.Background(), func() (*Server, error) { attempts++; return nil, want })
	if !errors.Is(err, want) || attempts != 1 {
		t.Fatalf("err=%v attempts=%d", err, attempts)
	}
}
func TestStartupLiveDuringOutage(t *testing.T) {
	h := &startupHandler{}
	r := httptest.NewRecorder()
	h.ServeHTTP(r, httptest.NewRequest("GET", "/api/live", nil))
	if r.Code != 200 {
		t.Fatal(r.Code)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if !errors.Is(h.initialize(ctx, func() (*Server, error) { t.Fatal("unexpected build"); return nil, nil }), context.Canceled) {
		t.Fatal("expected cancellation")
	}
}
