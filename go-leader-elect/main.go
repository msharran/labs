package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/urfave/cli"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Serving request", "url", r.URL.RawPath, "method", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("server is up"))
	})

	app := &cli.App{
		Name:  "server",
		Usage: "server runs demo raft leader election protocol",
		Action: func(ctx *cli.Context) error {
			slog.Info("Starting server")
			return http.ListenAndServe(":9900", mux)
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
