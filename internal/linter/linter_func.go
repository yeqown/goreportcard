package linter

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/yeqown/goreportcard/internal/types"
	"github.com/yeqown/log"
)

// Context is Lint function's param
type Context struct {
	Dir       string   // Dir of repo
	Filenames []string // Filenames of repo
	Branch    string   // Branch of repo
}

// Lint executes all checks on the given directory
//
// 1. get repo status: @fileCount @lineCount
// 2. call `golangci-lint` to lint, get errors
// 3. calc score of each linters
// 4. return result
func Lint(ctx Context) (result types.LintResult, err error) {
	log.Debugf("Lint recv params @dir=%s", ctx.Dir)

	filenames, err := visitGoFiles(ctx.Dir)
	if err != nil {
		err = errors.Errorf("could not get filenames: %v", err)
		return
	}
	if len(filenames) == 0 {
		err = errors.Errorf("no .go files found")
		return
	}
	ctx.Filenames = filenames

	var (
		linters   = getLinters()
		chanScore = make(chan types.Score, len(linters))
	)

	for _, linter := range linters {
		go execLinter(ctx, linter, chanScore)
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

	result = types.LintResult{
		Files:   len(filenames),
		Issues:  issuesCnt,
		Average: total,
		Scores:  scores,
		Grade:   types.GradeFromPercentage(total * 100),
	}

	return
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
func getLinters() []ILinter {
	return []ILinter{
		builtin{
			name: "govet", weight: .30,
			desc: "Vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string.",
		}, // govet
		builtin{
			name: "errcheck", weight: .10,
			desc: "Errcheck is a program for checking for unchecked errors in go programs. These unchecked errors can be critical bugs in some cases.",
		}, // errcheck
		builtin{
			name: "ineffassign", weight: .05,
			desc: "Detects when assignments to existing variables are not used.",
		}, // ineffassign
		builtin{
			name: "deadcode", weight: .05,
			desc: "Finds unused code",
		}, // deadcode
		builtin{
			name: "gosimple", weight: .05,
			desc: "Linter for Go source code that specializes in simplifying a code.",
		}, // gosimple
		builtin{
			name: "staticcheck", weight: .05,
			desc: "Staticcheck is a go vet on steroids, applying a ton of static analysis checks.",
		}, // staticcheck
		builtin{
			name: "structcheck", weight: .05,
			desc: "Finds unused struct fields.",
		}, // structcheck
		builtin{
			name: "unused", weight: .10,
			desc: "Scores Go code for unused constants, variables, functions and types.",
		}, // unused
		builtin{
			name: "varcheck", weight: .05,
			desc: "Finds unused global variables and constants.",
		}, // varcheck
		builtin{
			name: "typecheck", weight: .05,
			desc: "Like the front-end of a Go compiler, parses and type-checks Go codes.",
		}, // typecheck
		builtin{
			name: "funlen", weight: .10,
			desc: "Tool for detection of long functions.",
		}, // funlen
		builtin{
			name: "lll", weight: .10,
			desc: "Reports long lines.",
		}, // lll
		builtin{
			name: "nestif", weight: .15,
			desc: "Reports deeply nested if statements.",
		}, // nestif
	}
}

// execLinter exec linter.Execute and send types.Score by `chanScore`
func execLinter(ctx Context, linter ILinter, chanScore chan<- types.Score) {
	var errMsg string
	p, summaries, err := linter.Execute(ctx)
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
