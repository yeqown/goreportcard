package linter

import (
	"github.com/yeqown/goreportcard/internal/types"
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
	Percentage() (float64, []types.FileSummary, error)
}
