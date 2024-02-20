// Package server provides a server with routes.
package server

import (
	"context"
	"go-h1/internal/managers/secrets"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Flags uint8

type ServerOpts struct {
	Ctx   context.Context
	Flags Flags
	Addr  string
	Log   *slog.Logger
}

const (
	FLAG_DISABLE_LOGGING Flags = 1 << iota
	FLAG_DISABLE_ADMIN   Flags = 1 << iota
)

type Server struct {
	ctx     context.Context
	cancel  context.CancelFunc
	mux     *http.ServeMux
	httpsvr *http.Server
	log     *slog.Logger
	f       Flags
	sm      *secrets.SecretsManager
	wg      sync.WaitGroup
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
func NewServer(o ServerOpts) *Server {
	l := o.Log
	if o.Flags&FLAG_DISABLE_LOGGING != 0 {
		l = slog.New(&discardSlogHandler{})
	}

	ctx, cancel := context.WithCancel(o.Ctx)

	s := &Server{
		ctx:    ctx,
		cancel: cancel,
		log:    l,
		mux:    http.NewServeMux(),
		sm:     secrets.NewManager(ctx, l),
	}
	s.routes()

	s.httpsvr = &http.Server{
		Addr:    o.Addr,
		Handler: s, // s implements http.Handler plus easier to test with middleware
	}

	return s
}

func (s *Server) Run() error {
	s.log.Info("server started")
	defer s.log.Info("server stopped")

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		// stops sm when context is cancelled
		s.sm.Run()
	}()

	s.wg.Add(1)
	go func() {
		defer func() {
			s.wg.Done()
			s.log.Info("http server stopped")
		}()

		// stops server when context is cancelled
		<-s.ctx.Done()
		ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
		defer cancel()
		s.httpsvr.Shutdown(ctx)
	}()

	if err := s.httpsvr.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Wait() {
	s.cancel()
	s.wg.Wait()
	s.log.Info("stopped all goroutines managed by server")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	s.mux.ServeHTTP(w, r)
	s.log.Info(r.URL.Path, "method", r.Method, "duration", time.Since(start))
}

type discardSlogHandler struct{}

func (d *discardSlogHandler) Handle(context.Context, slog.Record) error { return nil }

func (d *discardSlogHandler) Enabled(context.Context, slog.Level) bool { return false }

func (d *discardSlogHandler) WithAttrs([]slog.Attr) slog.Handler { return d }

func (d *discardSlogHandler) WithGroup(string) slog.Handler { return d }
