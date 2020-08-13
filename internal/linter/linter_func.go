package linter

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/yeqown/goreportcard/internal/types"
	"github.com/yeqown/log"
)

// Lint executes all checks on the given directory
//
// 1. get repo status: @fileCount @lineCount
// 2. call `golangci-lint` to lint, get errors
// 3. calc score of each linters
// 4. return result
func Lint(dir string) (types.ChecksResult, error) {
	log.Debugf("Lint recv params @dir=%s", dir)

	filenames, err := visitGoFiles(dir)
	if err != nil {
		return types.ChecksResult{}, errors.Errorf("could not get filenames: %v", err)
	}
	if len(filenames) == 0 {
		return types.ChecksResult{}, errors.Errorf("no .go files found")
	}

	var (
		linters   = getLinters(dir, filenames)
		chanScore = make(chan types.Score)
	)

	for _, linter := range linters {
		go execLinter(linter, chanScore)
	}

	var (
		total, totalWeight float64
		issuesCnt          int
		n                  = len(linters)
		scores             = make(types.ByWeight, 0, 64)
	)

	// calc grade and score, then save into `types.CheckResult`
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

	r := types.ChecksResult{
		Files:   len(filenames),
		Issues:  issuesCnt,
		Average: total,
		Scores:  scores,
		Grade:   types.GradeFromPercentage(total * 100),
	}

	return r, nil
}

// https://golangci-lint.run/usage/linters/
//
// govet - Vet examines Go source code and reports suspicious constructs,
//         such as Printf calls whose arguments do not align with the format string
// errcheck - Errcheck is a program for checking for unchecked errors in go programs.
//         These unchecked errors can be critical bugs in some cases
// staticcheck - Staticcheck is a go vet on steroids, applying a ton of static analysis checks
// unused - Scores Go code for unused constants, variables, functions and types
// gosimple - Linter for Go source code that specializes in simplifying a code
// structcheck - Finds unused struct fields
// varcheck - Finds unused global variables and constants
// ineffassign - Detects when assignments to existing variables are not used
// deadcode - Finds unused code
// typecheck - Like the front-end of a Go compiler, parses and type-checks Go code
//

// getLinters . load all linters to run
// linters: https://golangci-lint.run/usage/linters/
func getLinters(dir string, filenames []string) []ILinter {
	return []ILinter{
		builtin{
			Dir: dir, Filenames: filenames, name: "govet", weight: .25,
			desc: "Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string.",
		}, // govet
		builtin{
			Dir: dir, Filenames: filenames, name: "errcheck", weight: .05,
			desc: "Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases.",
		}, // errcheck
		builtin{
			Dir: dir, Filenames: filenames, name: "ineffassign", weight: .05,
			desc: "Detects when assignments to existing variables are not used.",
		}, // ineffassign
		builtin{
			Dir: dir, Filenames: filenames, name: "deadcode", weight: .05,
			desc: "Finds unused code",
		}, // deadcode
		builtin{
			Dir: dir, Filenames: filenames, name: "gosimple", weight: .05,
			desc: "Linter for Go source code that specializes in simplifying a code.",
		}, // gosimple
		builtin{
			Dir: dir, Filenames: filenames, name: "staticcheck", weight: .05,
			desc: "Staticcheck is a go vet on steroids, applying a ton of static analysis checks.",
		}, // staticcheck
		builtin{
			Dir: dir, Filenames: filenames, name: "structcheck", weight: .05,
			desc: "Finds unused struct fields.",
		}, // structcheck
		builtin{
			Dir: dir, Filenames: filenames, name: "unused", weight: .05,
			desc: "Scores Go code for unused constants, variables, functions and types.",
		}, // unused
		builtin{
			Dir: dir, Filenames: filenames, name: "varcheck", weight: .05,
			desc: "Finds unused global variables and constants.",
		}, // varcheck
		builtin{
			Dir: dir, Filenames: filenames, name: "typecheck", weight: .05,
			desc: "Like the front-end of a Go compiler, parses and type-checks Go codes.",
		}, // typecheck
		builtin{
			Dir: dir, Filenames: filenames, name: "funlen", weight: .05,
			desc: "Tool for detection of long functions.",
		}, // funlen
		builtin{
			Dir: dir, Filenames: filenames, name: "lll", weight: .05,
			desc: "Reports long lines.",
		}, // lll
		builtin{
			Dir: dir, Filenames: filenames, name: "nestif", weight: .05,
			desc: "Reports deeply nested if statements.",
		}, // nestif
	}
}

// execLinter exec linter.Percentage and send types.Score by `chanScore`
func execLinter(linter ILinter, chanScore chan<- types.Score) {
	var errMsg string
	p, summaries, err := linter.Percentage()
	if err != nil {
		log.Errorf("Lint run linter=%s failed, err=%v", linter.Name(), err)
		errMsg = err.Error()
	}

	// send score to channel
	score := types.Score{
		Name:       linter.Name(),
		Desc:       linter.Description(),
		Summaries:  summaries,
		Weight:     linter.Weight(),
		Percentage: p,
		Error:      errMsg,
	}
	chanScore <- score
}
