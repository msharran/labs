package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"log/slog"
)

func main() {
	log := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.UnixDate,
		}),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log := log.WithGroup("request")
		
		log.Info("handling request", "url", r.Host)
		fmt.Fprintln(w, "Hello, World")
		log.Info("successfully handled request", "status", http.StatusOK)
	})

	log.Info("started the server", "addr", ":8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error("failed to handle requests", "err", err)
		os.Exit(1)
	}
}
