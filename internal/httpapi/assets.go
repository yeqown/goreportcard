package httpapi

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/yeqown/goreportcard/internal/types"
	"github.com/yeqown/log"
)

type assetsHandler struct {
	badgeCache sync.Map
}

// NewAssetsHandler to expose assets related handlers
func NewAssetsHandler() *assetsHandler {
	return &assetsHandler{
		badgeCache: sync.Map{},
	}
}

// AssetsHandler handles serving static files
func (hdl *assetsHandler) Assets(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "max-age=86400")
	assetURI := req.URL.Path
	http.ServeFile(w, req, assetURI[1:])
}

// FaviconHandler handles serving the favicon.ico
func (hdl *assetsHandler) Favicon(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "assets/favicon.ico")
}

// BadgeHandler handles fetching the badge images
// See: http://shields.io/#styles
func (hdl *assetsHandler) Badge(w http.ResponseWriter, req *http.Request, p *types.RepoReportParam) {
	style := req.URL.Query().Get("style")
	if style == "" {
		style = "flat"
	}

	if g, ok := hdl.badgeCache.Load(p.RepoIdentity()); ok {
		log.WithFields(log.Fields{
			"param":    p,
			"identity": p.RepoIdentity(),
		}).Infof("Fetching badge from cache")

		w.Header().Set("Cache-control", "no-store, no-badgeCache, must-revalidate")
		http.ServeFile(w, req, badgePath(g.(types.Grade), style))
		return
	}

	// not found in cache, then reload from lint
	r, err := doling(p, false)
	if err != nil {
		log.WithFields(log.Fields{
			"param": p,
			"error": err,
		}).Errorf("fetching badge failed")
		url := "https://img.shields.io/badge/go%20report-error-lightgrey.svg?style=" + style
		http.Redirect(w, req, url, http.StatusTemporaryRedirect)
		return
	}

	// update cache
	hdl.badgeCache.Store(p.RepoIdentity(), r.Grade)

	w.Header().Set("Cache-control", "no-store, no-badgeCache, must-revalidate")
	http.ServeFile(w, req, badgePath(r.Grade, style))
}

func badgePath(grade types.Grade, style string) string {
	return fmt.Sprintf("assets/badges/%s_%s.svg",
		strings.ToLower(string(grade)), strings.ToLower(style))
}
