package linter

import "github.com/yeqown/goreportcard/internal/types"

var _ ILinter = &builtin{}

type builtin struct {
	Dir       string
	Filenames []string

	name   string
	desc   string
	weight float64
}

func (b builtin) Name() string {
	return b.name
}

func (b builtin) Description() string {
	return b.desc
}

func (b builtin) Weight() float64 {
	return b.weight
}

func (b builtin) Percentage() (float64, []types.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=" + b.name,
		"--allow-parallel-runners",
		"--skip-dirs-use-default",
		"--tests=false",
	}

	return cmdHelper(b.Dir, b.Filenames, command)
}
