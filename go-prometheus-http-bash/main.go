package main

import (
	"encoding/json"
	"net/http"

	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricsNamespace = "cicd"

	stageRunTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "pipeline_stage_run_total",
			Help:      "",
		},
		[]string{"pipeline", "stage"},
	)

	stageRunDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "pipeline_stage_run_duration_seconds",
			Help:      "",
			Buckets:   prometheus.ExponentialBuckets(0.01, 2, 15),
		},
		[]string{"pipeline", "stage"},
	)
)

func main() {
	// We use a registry here to benefit from the consistency checks that
	// happen during registration.
	registry := prometheus.NewRegistry()
	registry.MustRegister(stageRunTotal, stageRunDurationSeconds)

	http.Handle("/metrics", promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	))

	http.HandleFunc("/metrics/pipeline/stage", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.String())
		// get pipeline_name, stage_name, duration_seconds from request body
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			PipelineName string  `json:"pipeline_name"`
			StageName    string  `json:"stage_name"`
			Status       string  `json:"status"`
			Duration     float64 `json:"duration_seconds"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Println("Error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("pipeline_name: %s, stage_name: %s, duration_seconds: %f\n", body.PipelineName, body.StageName, body.Duration)
		stageRunTotal.WithLabelValues(body.PipelineName, body.StageName).Inc()
		stageRunDurationSeconds.WithLabelValues(body.PipelineName, body.StageName).Observe(body.Duration)
		w.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe(":9900", nil)
	if err != nil {
		log.Fatal(err, "metrics endpoint listener failed")
	}
}
