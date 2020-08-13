package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/yeqown/goreportcard/internal/httpapi"
	"github.com/yeqown/goreportcard/internal/repository"
	"github.com/yeqown/goreportcard/internal/types"
	vcs "github.com/yeqown/goreportcard/internal/vcs-helper"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yeqown/log"
)

func startWebServer(cfg *types.Config) error {
	log.Debugf("start web with config=%+v", cfg)

	// load VCS downloader and others
	if err := vcs.Init(vcs.BuiltinTool, cfg.VCSOptions); err != nil {
		return errors.Wrap(err, "startWebServer.httpapi.Init")
	}

	// load db
	if err := repository.Init(cfg.DB); err != nil {
		return errors.Wrap(err, "startWebServer.repository.Init")
	}
	defer repository.GetRepo().Close()

	// prepare dir in which repos storage
	if err := os.MkdirAll(cfg.RepoRoot, 0755); err != nil && !os.IsExist(err) {
		return errors.Wrapf(err, "os mkdir in: %s", cfg.RepoRoot)
	}

	assetHdl := httpapi.NewAssetsHandler()
	http.HandleFunc("/", withMetrics(httpapi.HomeHandler))
	http.HandleFunc("/assets/", withMetrics(assetHdl.Assets))
	http.HandleFunc("/favicon.ico", withMetrics(assetHdl.Favicon))
	http.HandleFunc("/checks", withMetrics(httpapi.LintHandler))
	http.HandleFunc("/high_scores/", withMetrics(httpapi.HighScoresHandler))
	http.HandleFunc("/about/", withMetrics(httpapi.AboutHandler))
	http.HandleFunc("/report/", withMetrics(resolveRepoPath("report", httpapi.ReportHandler)))
	http.HandleFunc("/badge/", withMetrics(resolveRepoPath("badge", assetHdl.Badge)))

	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	log.Infof("Running on http://%s ...", addr)

	return http.ListenAndServe(addr, nil)
}

var (
	once                sync.Once
	responseTimeSummary *prometheus.SummaryVec
)

func withMetrics(fn http.HandlerFunc) http.HandlerFunc {
	once.Do(func() {
		responseTimeSummary = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "response_time_ms",
				Help: "Time to serve requests",
			}, // opts
			[]string{"endpoint"}, // label names
		)
		prometheus.MustRegister(responseTimeSummary)
	})

	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"query":  r.URL.RawQuery,
			"path":   r.URL.Path,
		}).Info("a request coming")

		start := time.Now()
		fn(w, r)
		elapsed := time.Since(start)
		responseTimeSummary.WithLabelValues(r.URL.Path).
			Observe(float64(elapsed.Nanoseconds()))
	}
}

type repoHandleFunc func(w http.ResponseWriter, req *http.Request, p *types.RepoReportParam)

// resolveRepoPath to resolve sub path of repo URL:
// https://example.com/report/github.com/golang/go
func resolveRepoPath(prefix string, fn repoHandleFunc) http.HandlerFunc {
	reg := regexp.MustCompile(fmt.Sprintf(`^/%s/([a-zA-Z0-9\-_\/\.]+)$`, prefix))

	return func(w http.ResponseWriter, r *http.Request) {
		branch := r.FormValue(types.BranchFormKey)
		if branch == "" {
			branch = types.MasterBranch
		}

		m := reg.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		if len(m) < 1 || m[1] == "" {
			http.Error(w, "Please enter a repository", http.StatusBadRequest)
			return
		}

		repo := m[1]

		// for backwards-compatibility, we must support URLs formatted as
		//   /report/[org]/[repo]
		// and they will be assumed to be github.com URLs. This is because
		// at first Go Report Card only supported github.com URLs, and
		// took only the org prefix and repo prefix as parameters. This is no longer the
		// case, but we do not want external links to break.
		oldFormat := regexp.MustCompile(fmt.Sprintf(`^/%s/([a-zA-Z0-9\-_]+)/([a-zA-Z0-9\-_]+)$`, prefix))
		m2 := oldFormat.FindStringSubmatch(r.URL.Path)
		if m2 != nil {
			// old format is being used
			repo = "github.com/" + repo
			log.Infof("Assuming intended repo is %q, redirecting", repo)
			http.Redirect(w, r, fmt.Sprintf("/%s/%s", prefix, repo), http.StatusMovedPermanently)
			return
		}

		fn(w, r, types.NewRepoParam(repo, branch))
	}
}
