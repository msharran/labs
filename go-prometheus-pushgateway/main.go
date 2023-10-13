package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	metricsNamespace = "caterpillar"

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
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10),
		},
		[]string{"pipeline", "stage"},
	)

	// successTime = prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "db_backup_last_success_timestamp_seconds",
	// 	Help: "The timestamp of the last successful completion of a DB backup.",
	// })
	// duration = prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "db_backup_duration_seconds",
	// 	Help: "The duration of the last DB backup in seconds.",
	// })
	// records = prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "db_backup_records_processed",
	// 	Help: "The number of records processed in the last DB backup.",
	// })
)

func main() {
	// We use a registry here to benefit from the consistency checks that
	// happen during registration.
	registry := prometheus.NewRegistry()
	registry.MustRegister(stageRunTotal, stageRunDurationSeconds)

	pusher := push.New("http://localhost:9091", "caterpillar").Gatherer(registry)

	buildStage()
	deployStage()
	if err := pusher.Push(); err != nil {
		fmt.Println("Could not push to Pushgateway:", err)
	}
}

func buildStage() {
	fmt.Println("building app")
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		stageRunDurationSeconds.WithLabelValues("foo-job-prod", "build").Observe(duration.Seconds())
	}()
	// wait for a random seconds to simulate build time
	max := 5
	min := 1
	time.Sleep(time.Duration(rand.Intn(max-min+1)+min) * time.Second)

	stageRunTotal.WithLabelValues("foo-job-prod", "build").Inc()
}

func deployStage() {
	fmt.Println("deploying app")
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		stageRunDurationSeconds.WithLabelValues("foo-job-prod", "deploy").Observe(duration.Seconds())
	}()
	max := 5
	min := 1
	time.Sleep(time.Duration(rand.Intn(max-min+1)+min) * time.Second)
	stageRunTotal.WithLabelValues("foo-job-prod", "deploy").Inc()
}
