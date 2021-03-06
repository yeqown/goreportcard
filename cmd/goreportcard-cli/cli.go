package main

import (
	"fmt"

	"github.com/yeqown/goreportcard/internal/types"

	"github.com/yeqown/goreportcard/internal/linter"

	"github.com/pkg/errors"
	"github.com/yeqown/log"
)

func runCli(dir string, verbose bool) error {
	log.SetLogLevel(log.LevelError)

	ctx := linter.Context{
		Dir:    dir,
		Branch: types.MasterBranch,
	}

	r, err := linter.Lint(ctx)
	if err != nil {
		log.Errorf("Fatal error checking %s: %s", dir, err.Error())
		return errors.Wrapf(err, "Fatal error checking: [%s]", dir)
	}

	fmt.Printf("Grade: %s (%.1f%%)\n", r.Grade, r.Average*100)
	fmt.Printf("FilesCount: %d\n", r.Files)
	fmt.Printf("IssuesCount: %d\n", r.Issues)

	for _, score := range r.Scores {
		fmt.Printf("%s: %d%%\n", score.Name, int64(score.Percentage*100))
		if verbose && len(score.Summaries) > 0 {
			for _, summary := range score.Summaries {
				fmt.Printf("\t%s\n", summary.Filename)
				for _, err := range summary.Errors {
					fmt.Printf("\t\tLine %d: %s\n", err.LineNumber, err.ErrorString)
				}
			}
		}
	}

	return nil
}
