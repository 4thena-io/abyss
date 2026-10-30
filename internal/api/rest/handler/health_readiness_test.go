package handler

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/4thena-io/abyss/internal/service"
)

func TestReadinessAndLiveness(t *testing.T) {
	online := false
	svc := service.NewHealthService("", "github-actions", "").WithDependencies(func(context.Context) error {
		if !online {
			return errors.New("down")
		}
		return nil
	}, nil)
	h := NewHealthHandler(svc)
	r := httptest.NewRecorder()
	h.Ready(r, httptest.NewRequest("GET", "/api/ready", nil))
	if r.Code != 503 {
		t.Fatal(r.Code)
	}
	r = httptest.NewRecorder()
	h.Live(r, httptest.NewRequest("GET", "/api/live", nil))
	if r.Code != 200 {
		t.Fatal(r.Code)
	}
	online = true
	r = httptest.NewRecorder()
	h.Ready(r, httptest.NewRequest("GET", "/api/ready", nil))
	if r.Code != 200 {
		t.Fatal(r.Code)
	}
}
