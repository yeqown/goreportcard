package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gojp/goreportcard/internal/httpapi"
	"github.com/gojp/goreportcard/internal/model"
	"github.com/gojp/goreportcard/internal/repository"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yeqown/log"
)

func startWebServer(port int) error {
	httpapi.Init()
	if err := repository.Init(); err != nil {
		return errors.Wrap(err, "startWebServer.repository.Init")
	}
	defer repository.GetRepo().Close()

	if err := os.MkdirAll(model.GetConfig().RepoRoot, 0755); err != nil && !os.IsExist(err) {
		log.Fatal("ERROR: could not create repos dir: ", err)
	}

	// prometheus metrics
	var m = newMetrics()

	http.HandleFunc(m.instrument("/assets/", httpapi.AssetsHandler))
	http.HandleFunc(m.instrument("/favicon.ico", httpapi.FaviconHandler))
	http.HandleFunc(m.instrument("/checks", httpapi.CheckHandler))
	http.HandleFunc(m.instrument("/report/", httpapi.MakeHandler("report", httpapi.ReportHandler)))
	http.HandleFunc(m.instrument("/badge/", httpapi.MakeHandler("badge", httpapi.BadgeHandler)))
	http.HandleFunc(m.instrument("/high_scores/", httpapi.HighScoresHandler))
	http.HandleFunc(m.instrument("/supporters/", httpapi.SupportersHandler))
	http.HandleFunc(m.instrument("/about/", httpapi.AboutHandler))
	http.HandleFunc(m.instrument("/", httpapi.HomeHandler))
	// register prometheus metrics handler
	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Infof("Running on http://%s ...", addr)

	return http.ListenAndServe(addr, nil)
}

// metrics provides functionality for monitoring the application statuks
type metrics struct {
	responseTimes *prometheus.SummaryVec
}

// newMetrics creates custom Prometheus metrics for monitoring
// application statistics.
func newMetrics() *metrics {
	m := &metrics{}
	m.responseTimes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "response_time_ms",
			Help: "Time to serve requests",
		}, // opts
		[]string{"endpoint"}, // label names
	)

	prometheus.MustRegister(m.responseTimes)
	return m
}

// recordDuration records the length of a request from start to now
func (m metrics) recordDuration(start time.Time, path string) {
	elapsed := time.Since(start)
	m.responseTimes.WithLabelValues(path).Observe(float64(elapsed.Nanoseconds()))
	log.Infof("[Request] Served %s in %s\n", path, elapsed)
}

// instrument adds metric instrumentation to the handler passed in as argument
func (m metrics) instrument(pattern string, h http.HandlerFunc) (string, http.HandlerFunc) {
	return pattern, func(w http.ResponseWriter, r *http.Request) {
		defer m.recordDuration(time.Now(), r.URL.Path)
		h.ServeHTTP(w, r)
	}
}
