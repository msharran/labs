package server

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type testdata struct {
	log *slog.Logger
	s   *Server
}

func setup(t *testing.T) (*testdata, func()) {
	t.Helper()

	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	s := NewServer(log)
	start := time.Now()

	return &testdata{log, s}, func() {
		log.Info("teardown", "name", t.Name(), "duration", time.Since(start))
	}
}

func TestAdminGet(t *testing.T) {
	t.Run("NotAdmin", func(t *testing.T) {
		td, teardown := setup(t)
		defer teardown()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/admin", nil)

		// call ServeHTTP to test with middleware
		td.s.ServeHTTP(w, r)

		if w.Code != http.StatusForbidden {
			t.Errorf("want status %d; got %d", http.StatusForbidden, w.Code)
		}

		if w.Body.String() != "not an admin\n" {
			t.Errorf("want body %q; got %q", "not an admin\n", w.Body.String())
		}
	})

	t.Run("Admin", func(t *testing.T) {
		td, teardown := setup(t)
		defer teardown()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/admin", nil)
		r.Header.Set("X-Admin", "true")

		// call ServeHTTP to test with middleware
		td.s.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("want status %d; got %d", http.StatusOK, w.Code)
		}

		if w.Body.String() != "Hello, Admin!" {
			t.Errorf("want body %q; got %q", "Hello, Admin!", w.Body.String())
		}
	})
}
