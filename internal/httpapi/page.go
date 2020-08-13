package httpapi

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/yeqown/goreportcard/internal/repository"
	"github.com/yeqown/goreportcard/internal/types"

	"github.com/dustin/go-humanize"
	"github.com/yeqown/log"
)

// AboutHandler handles the about page
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, http.StatusOK, tplAbout, nil)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		renderHTML(w, http.StatusNotFound, tpl404, nil)
	}
}

var cache struct {
	items []recentItemForDisplay
	mux   sync.Mutex
	count int
}

type recentItemForDisplay struct {
	Repo              string
	Grade             string
	Branch            string
	Score             string
	LastGeneratedTime string
}

// HomeHandler handles the homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] != "" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	var recentRepos []recentItemForDisplay

	cache.mux.Lock()
	cache.count++
	defer cache.mux.Unlock()

	if cache.count < 100 && len(cache.items) == 5 {
		recentRepos = cache.items
		log.Info("Fetching recent repos from cache...")
	} else {
		items, err := loadRecentlyViewed()
		if err != nil {
			log.Warnf("HomeHandler failed to loadRecentlyViewed, err=%v", err)
		}

		recentRepos = make([]recentItemForDisplay, len(items))
		var j = len(items) - 1
		for _, v := range items {
			recentRepos[j] = recentItemForDisplay{
				Repo:              v.Repo,
				Grade:             v.Grade,
				Branch:            v.Branch,
				Score:             fmt.Sprintf("%.2f", v.Score*100),
				LastGeneratedTime: humanize.Time(v.LastGeneratedTime),
			}
			j--
		}

		cache.items = recentRepos
		cache.count = 0
	}

	data := map[string]interface{}{
		"Recent": recentRepos,
	}
	renderHTML(w, http.StatusOK, tplHome, data)
}

// ReportHandler handles the report page
func ReportHandler(w http.ResponseWriter, r *http.Request, p *types.RepoReportParam) {
	log.WithFields(log.Fields{
		"param": p,
	}).Infof("displaying report")
	var loading bool

	lintResult, err := loadLintResult(p)
	if err != nil {
		switch err {
		case repository.ErrKeyNotFound:
			// don't bother logging - we already log in getFromCache. continue
		default:
			log.Errorf("ReportHandler failed load lintResult, err=%v", err)
		}
		loading = true
	}

	d, err := json.Marshal(lintResult)
	if err != nil {
		log.Errorf("ReportHandler: could not marshal JSON: err=%v", err)
		http.Error(w, "Failed to load cache object", 500)
		return
	}

	data := map[string]interface{}{
		"repo":     p.Repo(),
		"branch":   p.Branch(),
		"response": string(d),
		"loading":  loading,
		"domain":   types.GetConfig().Domain,
	}
	renderHTML(w, http.StatusOK, tplReport, data)
}

// HighScoresHandler handles the stats page
func HighScoresHandler(w http.ResponseWriter, r *http.Request) {
	var (
		reposCount int
		scores     ScoreHeap
		err        error
	)

	if scores, err = loadHighScores(); err != nil {
		log.Errorf("HighScoresHandler failed to loadHighScores, err=%v", err)
		Error(w, http.StatusInternalServerError, err)
		return
	}
	if reposCount, err = loadReposCount(); err != nil {
		log.Errorf("HighScoresHandler failed to loadReposCount, err=%v", err)
		Error(w, http.StatusInternalServerError, err)
		return
	}

	// handler scores
	sortedScores := make([]scoreItem, scores.Len())
	for i := range sortedScores {
		sortedScores[len(sortedScores)-i-1] = heap.Pop(&scores).(scoreItem)
	}

	data := map[string]interface{}{
		"HighScores": sortedScores,
		"Count":      reposCount,
	}
	renderHTML(w, http.StatusOK, tplHighscore, data)
}
