package linter

import "github.com/yeqown/goreportcard/internal/types"

var _ ILinter = &builtin{}

type builtin struct {
	name   string  // linter's name
	desc   string  // linter's desc
	weight float64 // linter's weight
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

func (b builtin) Execute(ctx Context) (float64, []types.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=" + b.name,
		"--allow-parallel-runners",
		"--skip-dirs-use-default=true",
		"--tests=false",
	}

	return cmdHelper(ctx, command)
}
