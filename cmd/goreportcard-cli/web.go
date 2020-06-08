package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gojp/goreportcard/internal/httpapi"
	"github.com/gojp/goreportcard/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yeqown/log"
)

// const (
// 	repoHome = ".repos/src"
// )

func startWebServer(port int) error {
	httpapi.Init()
	repository.Init() // TODO: handle error
	defer repository.GetRepo().Close()

	// TODO: 配置存放代码的位置可配置
	if err := os.MkdirAll("_repos/src", 0755); err != nil && !os.IsExist(err) {
		log.Fatal("ERROR: could not create repos dir: ", err)
	}

	// prometheus metrics
	var m = setupMetrics()

	http.HandleFunc(m.instrument("/assets/", httpapi.AssetsHandler))
	http.HandleFunc(m.instrument("/favicon.ico", httpapi.FaviconHandler))
	http.HandleFunc(m.instrument("/checks", httpapi.CheckHandler))
	http.HandleFunc(m.instrument("/report/", httpapi.MakeHandler("report", httpapi.ReportHandler)))
	http.HandleFunc(m.instrument("/badge/", httpapi.MakeHandler("badge", httpapi.BadgeHandler)))
	http.HandleFunc(m.instrument("/high_scores/", httpapi.HighScoresHandler))
	http.HandleFunc(m.instrument("/supporters/", httpapi.SupportersHandler))
	http.HandleFunc(m.instrument("/about/", httpapi.AboutHandler))
	http.HandleFunc(m.instrument("/", httpapi.HomeHandler))
	// register prometheus metrics hanlder
	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Infof("Running on http://%s ...", addr)

	return http.ListenAndServe(addr, nil)
}

// metrics provides functionality for monitoring the application statuks
type metrics struct {
	responseTimes *prometheus.SummaryVec
}

// setupMetrics creates custom Prometheus metrics for monitoring
// application statistics.
func setupMetrics() *metrics {
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
