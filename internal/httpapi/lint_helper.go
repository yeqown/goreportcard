package httpapi

import (
	"container/heap"
	"encoding/json"
	"strings"
	"time"

	"github.com/yeqown/goreportcard/internal/linter"
	"github.com/yeqown/goreportcard/internal/repository"
	"github.com/yeqown/goreportcard/internal/types"
	vcshelper "github.com/yeqown/goreportcard/internal/vcs-helper"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/yeqown/log"
)

// executing golangci-lint tool, and return result
func doling(p *types.RepoReportParam, forceRefresh bool) (types.LintResult, error) {
	log.WithFields(log.Fields{
		"param":        p,
		"forceRefresh": forceRefresh,
	}).Debugf("doling called")

	if !forceRefresh {
		resp, err := loadLintResult(p)
		if err == nil {
			return *resp, nil
		}
		// just log the error and continue
		log.Warnf("doling failed to loadLintResult, err=%v", err)
	}

	// fetch the repoIdentity and grade it
	root, err := vcshelper.GetDownloader().Download(p.Repo(), types.GetConfig().RepoRoot, p.Branch())
	if err != nil {
		return types.LintResult{}, errors.Errorf("could not clone repoIdentity: %v", err)
	}
	log.WithFields(log.Fields{
		"repo": p.Repo(),
	}).Infof("repo has been downloaded")

	result, err := linter.Lint(root)
	if err != nil {
		return types.LintResult{}, err
	}

	t := time.Now().UTC()
	lintResult := types.LintResult{
		Scores:               result.Scores,
		Average:              result.Average,
		Grade:                result.Grade,
		FilesCount:           result.Files,
		IssuesCount:          result.Issues,
		Repo:                 p.Repo(),
		ResolvedRepo:         p.Repo(),
		Branch:               p.Branch(),
		LastRefresh:          t,
		LastRefreshFormatted: t.Format(time.UnixDate),
		LastRefreshHumanized: humanize.Time(t),
	}

	var (
		isNewRepo bool // current repoIdentity is first encounter with goreportcard
		key       = lintResultKey(p)
	)

	v, err := repository.GetRepo().Get(key)
	if err != nil {
		log.Debugf("doling failed to getting lint result, key=%s, err=%v", key, err)
	}
	isNewRepo = len(v) == 0 || errors.Cause(err) == repository.ErrKeyNotFound

	// if this is a new repo, or the user force-refreshed, update the cache
	if isNewRepo || forceRefresh {
		if err = updateLintResult(p, lintResult); err != nil {
			log.Errorf("doling update lintResult failed key=%s, err=%v", key, err)
		}
		log.Debugf("doling updateLintResult success")
	}

	if err := updateMetadata(lintResult, p, isNewRepo); err != nil {
		log.Errorf("doling.updateMetadata failed: err=%v", err)
	}

	return lintResult, nil
}

// lintResultKey . to generate db.Key of lint result
func lintResultKey(p *types.RepoReportParam) []byte {
	return []byte("repos-" + p.RepoIdentity())
}

// loadLintResult query lintResult by repoIdentity, if hit in DB then return,
// otherwise return an error.
func loadLintResult(p *types.RepoReportParam) (*types.LintResult, error) {
	key := lintResultKey(p)
	data, err := repository.GetRepo().Get(key)
	if err != nil {
		return nil, err
	}

	// TRUE: hit cache
	resp := new(types.LintResult)
	if err = json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	resp.LastRefresh = resp.LastRefresh.UTC()
	resp.LastRefreshFormatted = resp.LastRefresh.Format(time.UnixDate)
	resp.LastRefreshHumanized = humanize.Time(resp.LastRefresh.UTC())
	resp.Grade = types.GradeFromPercentage(resp.Average * 100) // grade is not stored for some repos, yet
	return resp, nil
}

// updateLintResult update lintResult in DB.
func updateLintResult(p *types.RepoReportParam, result types.LintResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return errors.Wrap(err, "updateLintResult.jsonMarshal")
	}

	key := lintResultKey(p)
	if err = repository.GetRepo().Update(key, data); err != nil {
		return errors.Wrap(err, "updateLintResult.Update")
	}

	return nil
}

type recentItem struct {
	Repo              string
	Branch            string
	Grade             string
	Score             float64
	LastGeneratedTime time.Time
}

func (r recentItem) Equal(r2 recentItem) bool {
	return r.Repo == r2.Repo && r.Branch == r2.Branch
}

var (
	_recentKey   = []byte("recent")
	_scoreKey    = []byte("scores")
	_reposCntKey = []byte("total_repos")
)

// updateRecentlyViewed .
func updateRecentlyViewed(item recentItem) (err error) {
	var (
		items []recentItem
		_repo = repository.GetRepo()
		d     []byte
	)

	if d, err = _repo.Get(_recentKey); err != nil &&
		errors.Cause(err) != repository.ErrKeyNotFound {
		return errors.Wrap(err, "updateRecentlyViewed.repo.Get")
	}

	if len(d) != 0 {
		err = json.Unmarshal(d, &items)
		if err != nil {
			return errors.Wrap(err, "updateRecentlyViewed.jsonUnmarshal")
		}
	}

	// add it to the slice, if it is not in there already
	for _, v := range items {
		if v.Equal(item) {
			log.Infof("updateRecentlyViewed has exists repoIdentity=%s, so skipped", item.Repo)
			return
		}
	}

	items = append(items, item)
	if len(items) > 5 {
		// trim recent if it's grown to over 5
		items = (items)[1:6]
	}
	d, err = json.Marshal(&items)
	if err != nil {
		return errors.Wrap(err, "updateRecentlyViewed.jsonMarshal")
	}

	log.Debugf("updateRecentlyViewed will save key=%s, v=%s", _recentKey, d)
	return _repo.Update(_recentKey, d)
}

// loadRecentlyViewed .
func loadRecentlyViewed() ([]recentItem, error) {
	var (
		items []recentItem
		_repo = repository.GetRepo()
	)

	d, err := _repo.Get(_recentKey)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"d": string(d),
	}).Debugf("loadRecentlyViewed raw data")
	if err = json.Unmarshal(d, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// loadHighScores .
func loadHighScores() (scores ScoreHeap, err error) {
	var (
		_repo = repository.GetRepo()
		d     []byte
	)

	d, err = _repo.Get(_scoreKey)
	if err != nil {
		// if there's nothing to show, make it default empty
		if errors.Cause(err) == repository.ErrKeyNotFound {
			return scores, nil
		}
		err = errors.Wrap(err, "loadHighScores.repoGet")

		return
	}

	if err = json.Unmarshal(d, &scores); err != nil {
		err = errors.Wrap(err, "loadHighScores.jsonUnmarshal")
		return
	}
	return
}

// updateHighScores .
func updateHighScores(result types.LintResult, p *types.RepoReportParam) (err error) {
	var (
		_repo = repository.GetRepo()
		d     []byte
	)

	// check if we need to update the high score list
	// only repos with 100+ files are considered for the high score list
	// TODO: make this as configable
	if result.FilesCount < 10 {
		return nil
	}

	if d, err = _repo.Get(_scoreKey); err != nil &&
		errors.Cause(err) != repository.ErrKeyNotFound {

		return
	}

	var scores = new(ScoreHeap)
	if len(d) != 0 {
		if err = json.Unmarshal(d, scores); err != nil {
			return err
		}
	}

	if len(*scores) > 0 && (*scores)[0].Score > result.Average*100.0 && len(*scores) == 50 {
		// lowest score on list is higher than this repo's score, so no need to add, unless
		// we do not have 50 high scores yet
		return nil
	}

	// if this repo is already in the list, remove the original entry:
	for idx, v := range *scores {
		if strings.EqualFold(v.Repo, p.RepoIdentity()) {
			heap.Remove(scores, idx)
			break
		}
	}

	// now we can safely push it onto the heap
	heap.Init(scores)
	heap.Push(scores, scoreItem{
		Repo:   p.Repo(),
		Branch: p.Branch(),
		Score:  result.Average * 100.0,
		Files:  result.FilesCount,
	})

	if len(*scores) > 50 {
		// trim heap if it's grown to over 50
		*scores = (*scores)[1:51]
	}

	// save back
	if d, err = json.Marshal(&scores); err != nil {
		return err
	}

	return _repo.Update(_scoreKey, d)
}

// loadReposCount .
func loadReposCount() (cnt int, err error) {
	var (
		_repo = repository.GetRepo()
		d     []byte
	)

	d, err = _repo.Get(_reposCntKey)
	if err != nil {
		if errors.Cause(err) == repository.ErrKeyNotFound {
			return cnt, nil
		}
		err = errors.Wrap(err, "loadReposCount.repoGet")

		return
	}

	if err = json.Unmarshal(d, &cnt); err != nil {

		return
	}
	return
}

// only new can update
func incrReposCnt(repoIdentity string) (err error) {
	log.Infof("New repo %q, adding to repo count...", repoIdentity)
	var (
		_repo = repository.GetRepo()
		d     []byte
		cnt   int
	)

	// load and unmarshal
	if d, err = _repo.Get(_reposCntKey); err != nil &&
		errors.Cause(err) != repository.ErrKeyNotFound {

		return err
	}
	if len(d) != 0 {
		if err := json.Unmarshal(d, &cnt); err != nil {
			return err
		}
	}

	cnt++
	if d, err = json.Marshal(cnt); err != nil {
		return err
	}
	if err = _repo.Update(_reposCntKey, d); err != nil {
		return err
	}

	return nil
}

// updateMetadata to record some data of goreportcard
func updateMetadata(result types.LintResult, p *types.RepoReportParam, isNewRepo bool) (err error) {
	// increase repos count
	if isNewRepo {
		if err = incrReposCnt(p.RepoIdentity()); err != nil {
			log.Errorf("updateMetadata.incrReposCnt failed: err=%v", err)
		}
	}

	item := recentItem{
		Repo:              p.Repo(),
		Branch:            p.Branch(),
		Grade:             string(result.Grade),
		Score:             result.Average,
		LastGeneratedTime: result.LastRefresh,
	}

	if err = updateRecentlyViewed(item); err != nil {
		log.Errorf("updateMetadata.updateRecentlyViewed failed: err=%v", err)
	}
	if err = updateHighScores(result, p); err != nil {
		log.Errorf("updateMetadata.updateHighScores failed: err=%v", err)
	}

	log.Infof("updateMetadata success")
	return
}
