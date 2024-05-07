package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/lmittmann/tint"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	cf     *Conf
	e      *echo.Echo
	log    *slog.Logger
	wg     *sync.WaitGroup
}

func New() *App {
	// Load env vars from .env file (for local development),
	// if env vars are already set, they will not be overridden
	_ = godotenv.Load()

	a := &App{
		e:  echo.New(),
		wg: &sync.WaitGroup{},
	}
	cf, err := newConfigFromEnv()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	a.cf = cf
	a.log = newLoggerFromConfig(cf)
	a.ctx, a.cancel = a.newContextWithSignal()
	a.setupRoutes()
	return a
}

func (a *App) Ctx() context.Context {
	return a.ctx
}

func (a *App) Run() {
	a.wg.Add(1)
	go func() {
		defer func() {
			a.wg.Done()
			a.Logger().Info("server stopped")
		}()

		if err := a.e.Start(a.cf.Addr); err != nil && err != http.ErrServerClosed {
			a.Logger().Error("server error", "error", err)
			// cancel the context to signal Start() to return
			a.cancel()
		}
	}()

	a.Logger().Info("started server", "addr", a.cf.Addr)
	<-a.ctx.Done()
}

func (a *App) Shutdown(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := a.e.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		a.Logger().Error("failed to gracefully shutdown server", "error", err)
		return
	}

	a.Logger().Info("initiated graceful shutdown, waiting for goroutines to finish")
	a.wg.Wait()
	a.Logger().Info("server shutdown gracefully")
}

func (a *App) Logger() *slog.Logger {
	return a.log
}

func (a *App) setupRoutes() {
	// health check
	a.e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "PONG")
	})
}

func (a *App) newContextWithSignal() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-c:
			a.Logger().Info("signal received, cancelling context", "signal", sig.String())
			cancel()
		case <-ctx.Done():
		}
	}()

	cancelFunc := func() {
		signal.Stop(c)
		cancel()
	}

	return ctx, cancelFunc
}

func newLoggerFromConfig(cf *Conf) (l *slog.Logger) {
	w := os.Stdout

	lvl := slog.LevelInfo
	if cf.Verbose {
		lvl = slog.LevelDebug
	}

	var h slog.Handler
	switch cf.LogFormat {
	case "development":
		h = tint.NewHandler(w, &tint.Options{
			Level:      lvl,
			TimeFormat: time.Kitchen,
		})
	case "production":
		h = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: lvl,
		})
	default:
		panic(fmt.Errorf("invalid log format: %s", cf.LogFormat))
	}

	l = slog.New(h)
	l.Debug("logger created", "verbose", cf.Verbose, "log_format", cf.LogFormat)
	return
}
