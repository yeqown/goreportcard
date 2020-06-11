package httpapi

import (
	"errors"
	"net/http"
	"strings"
)

type repoPatternHandlerFn func(w http.ResponseWriter, req *http.Request, repo string)

// PathPatternHandler @uri=path/github.com/owner/repoName
// get @repo=github.com/owner/repoName then pass to `fn`
func PathPatternHandler(pathName string, fn repoPatternHandlerFn) http.HandlerFunc {
	purePathName := strings.TrimPrefix(pathName, "/")

	return func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.Path
		repoName, ok := getRepoFromURI(purePathName, uri)
		if !ok {
			Error(w, http.StatusBadRequest, errors.New("invalid repo path"))
			return
		}
		fn(w, r, repoName)
	}
}

func getRepoFromURI(pathName, uri string) (repoName string, ok bool) {
	prefix := "/" + pathName + "/"
	if !strings.HasPrefix(uri, prefix) {
		return
	}

	// get repoName from after path
	repoName = strings.TrimPrefix(uri, prefix)

	// to validate repoName
	if strings.Count(repoName, "/") != 2 {
		return
	}

	return repoName, true
}
