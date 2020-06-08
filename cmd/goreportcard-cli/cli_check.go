package main

import (
	"fmt"

	"github.com/gojp/goreportcard/internal/linter"

	"github.com/yeqown/log"
)

func cliCheck(dir string, verbose bool) error {
	log.SetLogLevel(log.LevelError)

	r, err := linter.Lint(dir)
	if err != nil {
		log.Errorf("Fatal error checking %s: %s", dir, err.Error())
		return err
	}

	fmt.Printf("Grade: %s (%.1f%%)\n", r.Grade, r.Average*100)
	fmt.Printf("Files: %d\n", r.Files)
	fmt.Printf("Issues: %d\n", r.Issues)

	for _, score := range r.Scores {
		fmt.Printf("%s: %d%%\n", score.Name, int64(score.Percentage*100))
		if verbose && len(score.FileSummaries) > 0 {
			for _, summary := range score.FileSummaries {
				fmt.Printf("\t%s\n", summary.Filename)
				for _, err := range summary.Errors {
					fmt.Printf("\t\tLine %d: %s\n", err.LineNumber, err.ErrorString)
				}
			}
		}
	}

	return nil
}
