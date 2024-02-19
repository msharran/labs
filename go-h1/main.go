package main

import (
	"context"
	"flag"
	"fmt"
	"go-h1/internal/server"
	"log/slog"
	"net/http"
	"os"
)

var (
	port           string
	disableLogging bool
	disableAdmin   bool
)

func init() {
	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.StringVar(&port, "p", "8080", "port to listen on (shorthand)")
}

func main() {
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	p, ok := os.LookupEnv("PORT")
	if ok {
		log.Info("PORT environment variable found", "port", p)
		port = p
	}

	ctx := server.WithLogger(context.Background(), log)
	s := server.NewServer(ctx, server.ServerOpts{
		// Flags: server.FLAG_DISABLE_LOGGING | server.FLAG_DISABLE_ADMIN,
	})

	log.Info("Starting server", "port", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), s)
	if err != nil {
		log.Error("Server failed", "error", err)
		os.Exit(1)
	}

	log.Info("Server stopped")
}
