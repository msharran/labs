package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/lmittmann/tint"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var duration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "The duration of requests processed by the server",
		Buckets: []float64{0.1, 0.3, 1.2, 5.0},
	},
	[]string{"path"},
)

func main() {
	log := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.UnixDate,
		}),
	)

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log := log.WithGroup("request")

		action := r.URL.Query().Get("action")
		if action == "start" {
			log.Info("starting request", "url", r.Host)
			duration.WithLabelValues(r.URL.Path).Observe(0)
		}

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
