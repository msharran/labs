package main

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	flag "github.com/spf13/pflag"
)

var (
	logLevel = &slog.LevelVar{}
	verbose  = false
)

func init() {
	flag.BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

func main() {
	flag.Parse()
	if verbose {
		logLevel.Set(slog.LevelDebug)
	}

	// slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	// slog.SetDefault(slogger)
	// slog.Info("Hello Gopher!",
	// 	"from", "slog",
	// 	"type", "text")
	// slog.Debug("Verbose log", "from", "slog")
	// slog.Error("Some fatal error", "from", "slog")

	// colored

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.UnixDate,
		}),
	))
	slog.Info("Hello Gopher!",
		"from", "slog",
		"type", "text")
	slog.Debug("Debug log", "from", "slog")
	slog.Warn("Warn log", "from", "slog")
	slog.Error("Some fatal error", "from", "slog", "err", errors.New("connection failed"))
}
