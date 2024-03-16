package main

import (
	"net/http"
	"os"
	"time"

	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tmpFilesCollectionDuraionMetrics = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: tmpStatsNamespace("scrape_duration_seconds"),
		Help: "Duration of /tmp file collection",
		// time in seconds bucket, starts with very small duration
		// and goes up to 5 seconds
		Buckets: prometheus.ExponentialBuckets(0.000001, 5, 15),
	})
)

func tmpStatsNamespace(s string) string {
	return "tmp_stats_" + s
}

type fileStat struct {
	name  string
	size  int64
	isDir bool
}

// fileCollector is inspired by
// baseGoCollector in client_golang/prometheus pkg
type fileCollector struct {
	// for constants use desc
	tmpFilesTotalDesc *prometheus.Desc

	tmpFilesDirTotal *prometheus.Desc

	tmpFileSizeMetrics *prometheus.Desc
}

// You must create a constructor for you collector that
// initializes every descriptor and returns a pointer to the collector
func newFileCollector() *fileCollector {
	return &fileCollector{
		tmpFilesTotalDesc: prometheus.NewDesc(
			tmpStatsNamespace("files_total"),
			"The total number of files in /tmp directory",
			nil,
			nil),
		tmpFilesDirTotal: prometheus.NewDesc(
			tmpStatsNamespace("directories_total"),
			"The total number of directories in /tmp directory",
			nil,
			nil),
		tmpFileSizeMetrics: prometheus.NewDesc(
			tmpStatsNamespace("file_size_bytes"),
			"The size of the file in /tmp directory",
			[]string{"file"},
			nil),
	}
}

// Describe returns all descriptions of the collector.
func (c *fileCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.tmpFilesTotalDesc
	ch <- c.tmpFilesDirTotal
	ch <- c.tmpFileSizeMetrics
}

// Collect returns the current state of all metrics of the collector.
func (c *fileCollector) Collect(ch chan<- prometheus.Metric) {
	// get the data from your store, here it
	// is the file system
	tempfiles := getTempFiles()

	ch <- prometheus.MustNewConstMetric(c.tmpFilesTotalDesc, prometheus.GaugeValue, float64(len(tempfiles)))

	var dirs int
	start := time.Now()
	for _, f := range tempfiles {
		if f.isDir {
			dirs++
		}
		// histograms can't be created with MustNewConstMetric
		// so we have to use the already created histogram
		// and add the data to it during collection
		ch <- prometheus.MustNewConstMetric(c.tmpFileSizeMetrics, prometheus.GaugeValue, float64(f.size), f.name)
	}
	end := time.Since(start)
	tmpFilesCollectionDuraionMetrics.Observe(end.Seconds())

	ch <- prometheus.MustNewConstMetric(c.tmpFilesDirTotal, prometheus.GaugeValue, float64(dirs))
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request %s %s", r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func main() {
	registry := prometheus.NewRegistry()
	registry.MustRegister(tmpFilesCollectionDuraionMetrics, newFileCollector())

	http.Handle("GET /metrics", loggingMiddleware(promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	)))

	log.Println("listening on :8080, serving metrics endpoint at /metrics")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err, "metrics endpoint listener failed")
	}
}

func getTempFiles() []*fileStat {
	files, err := os.ReadDir("/tmp/")
	if err != nil {
		log.Fatal(err)
	}

	var stats []*fileStat
	for _, file := range files {
		var size int64

		info, err := file.Info()
		if err != nil {
			log.Printf("error getting file info for %s: %v", file.Name(), err)
		} else {
			size = info.Size()
		}

		stats = append(stats, &fileStat{name: file.Name(), size: size, isDir: file.IsDir()})
	}
	return stats
}
