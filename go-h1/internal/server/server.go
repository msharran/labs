// Package server provides a server with routes.
package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Flags uint8

type ServerOpts struct {
	Flags Flags
}

const (
	FLAG_DISABLE_LOGGING Flags = 1 << iota
	FLAG_DISABLE_ADMIN   Flags = 1 << iota
)

type Server struct {
	ctx context.Context
	mux *http.ServeMux
	log *slog.Logger
	f   Flags
}

// NewServer creates a new server with a logger
// and a ServeMux with routes.
// Don't add dependencies to the constructor
// since tests will call this. Tests can set only
// required dependencies after calling NewServer.
//
// Example test:
//
//	s := server.NewServer(log)
//	s.db = mockDB
//	w := httptest.NewRecorder()
//	r := httptest.NewRequest("GET", "/path", nil)
//	s.ServeHTTP(w, r)
//	if w.Code != http.StatusOK {
//	  t.Errorf("want status %d; got %d", http.StatusOK, w.Code)
//	}
func NewServer(ctx context.Context, o ServerOpts) *Server {
	l := FromContext(ctx)

	if o.Flags&FLAG_DISABLE_LOGGING != 0 {
		l = slog.New(&discardSlogHandler{})
	}

	s := &Server{
		mux: http.NewServeMux(),
		log: l,
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	s.mux.ServeHTTP(w, r)
	s.log.Info(r.URL.Path, "method", r.Method, "duration", time.Since(start))
}

type logContextKey struct{}

func FromContext(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value(logContextKey{}).(*slog.Logger)
	if !ok {
		return slog.New(&discardSlogHandler{})
	}

	return l
}

func WithLogger(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, logContextKey{}, log)
}

type discardSlogHandler struct{}

func (d *discardSlogHandler) Handle(context.Context, slog.Record) error { return nil }

func (d *discardSlogHandler) Enabled(context.Context, slog.Level) bool { return false }

func (d *discardSlogHandler) WithAttrs([]slog.Attr) slog.Handler { return d }

func (d *discardSlogHandler) WithGroup(string) slog.Handler { return d }
