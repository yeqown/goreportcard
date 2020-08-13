package httpapi

import (
	"fmt"
	"net/http"

	"github.com/yeqown/goreportcard/internal/types"

	"github.com/pkg/errors"
	"github.com/yeqown/log"
)

const (
	_masterBranch  = "master"
	_repoFormKey   = "repo"
	_branchFormKey = "branch"
)

// LintHandler handles the request for checking a repo
func LintHandler(w http.ResponseWriter, r *http.Request) {
	repo := r.FormValue(_repoFormKey)
	branch := r.FormValue(_branchFormKey)

	// TODO: valid repo format "github.com/xxx/xxx"
	log.WithFields(log.Fields{
		"repo":   repo,
		"branch": branch,
	}).Infof("checking repo")

	if branch == "" {
		branch = _masterBranch
	}

	// if this is a GET request, try to fetch from cached version in badger first
	forceRefresh := r.Method != "GET"
	p := types.NewRepoParam(repo, branch)

	_, err := doling(p, forceRefresh)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorf("doling failed")

		Error(w, http.StatusBadRequest, errors.Wrap(err, "Could not analyze the repository"))
		return
	}

	data := map[string]string{
		"redirect": reportPageURI(repo, branch),
	}
	JSON(w, http.StatusOK, data)
}

func reportPageURI(repo, branch string) string {
	return fmt.Sprintf("/report/%s?branch=%s", repo, branch)
}
