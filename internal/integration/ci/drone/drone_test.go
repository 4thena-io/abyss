package drone

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDroneCi_DeleteRepo(t *testing.T) {
	for _, status := range []int{http.StatusNoContent, http.StatusInternalServerError} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete || r.URL.Path != "/api/repos/acme/my-app" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				if got := r.URL.Query().Get("remove"); got != "true" {
					t.Errorf("expected remove=true, got %q", got)
				}
				if got := r.Header.Get("Authorization"); got != "Bearer token" {
					t.Errorf("unexpected authorization header: %q", got)
				}
				w.WriteHeader(status)
			}))
			defer srv.Close()

			ci, err := NewDroneCi(srv.URL, "token")
			if err != nil {
				t.Fatalf("failed to create ci client: %v", err)
			}
			err = ci.DeleteRepo(context.Background(), 7, "acme/my-app")
			if (err != nil) != (status >= 400) {
				t.Fatalf("unexpected deletion result for status %d: %v", status, err)
			}
		})
	}
}
