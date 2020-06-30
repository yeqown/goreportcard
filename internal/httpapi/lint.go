package httpapi

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/yeqown/log"
)

const (
	// RepoPrefix is the badger prefix for repos
	RepoPrefix string = "repos-"
)

// LintHandler handles the request for checking a repo
func LintHandler(w http.ResponseWriter, r *http.Request) {
	repo := r.FormValue("repo")
	// TODO: valid repo format "github.com/xxx/xxx"
	log.Infof("Checking repo %q...", repo)

	// if this is a GET request, try to fetch from cached version in badger first
	forceRefresh := r.Method != "GET"

	_, err := doling(repo, forceRefresh)
	if err != nil {
		log.Errorf("doling failed, err=%v", err)
		Error(w, http.StatusBadRequest, errors.Wrap(err, "Could not analyze the repository"))
		return
	}

	JSON(w, http.StatusOK, map[string]string{"redirect": "/report/" + repo})
}
