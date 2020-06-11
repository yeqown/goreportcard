package linter

import (
	"fmt"
	"sort"

	"github.com/gojp/goreportcard/internal/model"

	"github.com/yeqown/log"
)

// ILinter describes what methods various checks (gofmt, go lint, etc.)
// should implement
type ILinter interface {
	// Name of ILinter
	Name() string

	// Desc of ILinter
	Description() string

	// Weight of ILinter to calc score
	Weight() float64

	// Percentage returns the passing percentage of the check,
	// as well as a map of filename to output
	Percentage() (float64, []model.FileSummary, error)
}

// Lint executes all checks on the given directory
//
// 1. get repo status: @fileCount @lineCount
// 2. call `golangci-lint` to lint, get errors
// 3. calc score of each linters
// 4. return result
func Lint(dir string) (model.ChecksResult, error) {
	log.Debugf("Lint recv params @dir=%s", dir)

	filenames, skipped, err := visitGoFiles(dir)
	if err != nil {
		return model.ChecksResult{}, fmt.Errorf("could not get filenames: %v", err)
	}
	_ = skipped
	if len(filenames) == 0 {
		return model.ChecksResult{}, fmt.Errorf("no .go files found")
	}

	var (
		linters   = getLinters(dir, filenames)
		chanScore = make(chan model.Score)
	)

	for _, linter := range linters {
		go execLinter(linter, chanScore)
	}

	var (
		total, totalWeight float64
		issuesCnt          int
		n                  = len(linters)
		scores             = make(model.ByWeight, 0, 64)
	)

	// calc grade and score, then save into `model.CheckResult`
	for i := 0; i < n; i++ {
		score := <-chanScore
		scores = append(scores, score)

		total += score.Percentage * score.Weight
		totalWeight += score.Weight
		for _, summary := range score.Summaries {
			issuesCnt += len(summary.Errors)
		}
	}
	close(chanScore)
	total /= totalWeight
	sort.Sort(scores)

	r := model.ChecksResult{
		Files:   len(filenames),
		Issues:  issuesCnt,
		Average: total,
		Scores:  scores,
		Grade:   model.GradeFromPercentage(total * 100),
	}

	return r, nil
}

// execLinter exec linter.Percentage and send model.Score by `chanScore`
func execLinter(linter ILinter, chanScore chan<- model.Score) {
	var errMsg string
	p, summaries, err := linter.Percentage()
	if err != nil {
		log.Errorf("Lint run linter=%s failed, err=%v", linter.Name(), err)
		errMsg = err.Error()
	}

	// send score to channel
	score := model.Score{
		Name:       linter.Name(),
		Desc:       linter.Description(),
		Summaries:  summaries,
		Weight:     linter.Weight(),
		Percentage: p,
		Error:      errMsg,
	}
	chanScore <- score
}
