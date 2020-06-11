package linter

import "github.com/gojp/goreportcard/internal/model"

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
		govet{Dir: dir, Filenames: filenames},       // govet
		errcheck{Dir: dir, Filenames: filenames},    // errcheck, disable errcheck for now, too slow and not finalized
		ineffassign{Dir: dir, Filenames: filenames}, // ineffassign
		deadcode{Dir: dir, Filenames: filenames},    // deadcode
		gosimple{Dir: dir, Filenames: filenames},    // gosimple
		staticcheck{Dir: dir, Filenames: filenames}, // staticcheck
		structcheck{Dir: dir, Filenames: filenames}, // structcheck
		unused{Dir: dir, Filenames: filenames},      // unused
		varcheck{Dir: dir, Filenames: filenames},    // varcheck
		typecheck{Dir: dir, Filenames: filenames},   // typecheck
	}
}

// govet is the check for the go vet command
type govet struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g govet) Name() string { return "govet" }

// Weight returns the weight this check has in the overall average
func (g govet) Weight() float64 { return .25 }

// Percentage returns the percentage of .go files that pass go vet
func (g govet) Percentage() (float64, []model.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=govet",
		"--allow-parallel-runners",
	}
	return cmdHelper(g.Dir, g.Filenames, command)
}

// Desc returns the description of go lint
func (g govet) Description() string {
	return `<code>go vet</code> examines Go source code and reports suspicious constructs, 
	such as Printf calls whose arguments do not align with the format string.`
}

// errcheck is the check for the errcheck command
type errcheck struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (c errcheck) Name() string { return "errcheck" }

// Weight returns the weight this check has in the overall average
func (c errcheck) Weight() float64 { return .15 }

// Percentage returns the percentage of .go files that pass gofmt
func (c errcheck) Percentage() (float64, []model.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=errcheck",
		"--allow-parallel-runners",
	}
	return cmdHelper(c.Dir, c.Filenames, command)
}

// Desc returns the description of gofmt
func (c errcheck) Description() string {
	return `<a href="https://github.com/kisielk/errcheck">errcheck</a> finds unchecked errors in go programs`
}

// ineffassign is the check for the ineffassign command
type ineffassign struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g ineffassign) Name() string { return "ineffassign" }

// Weight returns the weight this check has in the overall average
func (g ineffassign) Weight() float64 { return 0.05 }

// Percentage returns the percentage of .go files that pass gofmt
// golangci-lint run --deadline=180s --disable-all --enable=ineffassign
func (g ineffassign) Percentage() (float64, []model.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=ineffassign",
		"--allow-parallel-runners",
	}
	return cmdHelper(g.Dir, g.Filenames, command)
}

// Desc returns the description of ineffassign
func (g ineffassign) Description() string {
	return `<a href="https://github.com/gordonklaus/ineffassign">ineffassign</a> detects ineffectual assignments in Go code.`
}

// deadcode is the check for the go vet command
type deadcode struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g deadcode) Name() string { return "deadcode" }

// Weight returns the weight this check has in the overall average
func (g deadcode) Weight() float64 { return .25 }

// Percentage returns the percentage of .go files that pass go vet
func (g deadcode) Percentage() (float64, []model.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=deadcode",
		"--allow-parallel-runners",
	}
	return cmdHelper(g.Dir, g.Filenames, command)
}

// Desc returns the description of go lint
func (g deadcode) Description() string {
	return `<code>go vet</code> examines Go source code and reports suspicious constructs, 
	such as Printf calls whose arguments do not align with the format string.`
}

// gosimple is the check for the go vet command
//
type gosimple struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g gosimple) Name() string { return "gosimple" }

// Weight returns the weight this check has in the overall average
func (g gosimple) Weight() float64 { return .05 }

// Percentage returns the percentage of .go files that pass go vet
func (g gosimple) Percentage() (float64, []model.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=gosimple",
		"--allow-parallel-runners",
	}
	return cmdHelper(g.Dir, g.Filenames, command)
}

// Desc returns the description of go lint
func (g gosimple) Description() string {
	return `<code>go vet</code> examines Go source code and reports suspicious constructs, 
	such as Printf calls whose arguments do not align with the format string.`
}

// staticcheck is the check for the go vet command
//
type staticcheck struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g staticcheck) Name() string { return "staticcheck" }

// Weight returns the weight this check has in the overall average
func (g staticcheck) Weight() float64 { return .05 }

// Percentage returns the percentage of .go files that pass go vet
func (g staticcheck) Percentage() (float64, []model.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=staticcheck",
		"--allow-parallel-runners",
	}
	return cmdHelper(g.Dir, g.Filenames, command)
}

// Desc returns the description of go lint
func (g staticcheck) Description() string {
	return `<code>go vet</code> examines Go source code and reports suspicious constructs, 
	such as Printf calls whose arguments do not align with the format string.`
}

// structcheck is the check for the go vet command
//
type structcheck struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g structcheck) Name() string { return "structcheck" }

// Weight returns the weight this check has in the overall average
func (g structcheck) Weight() float64 { return .05 }

// Percentage returns the percentage of .go files that pass go vet
func (g structcheck) Percentage() (float64, []model.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=structcheck",
		"--allow-parallel-runners",
	}
	return cmdHelper(g.Dir, g.Filenames, command)
}

// Desc returns the description of go lint
func (g structcheck) Description() string {
	return `<code>go vet</code> examines Go source code and reports suspicious constructs, 
	such as Printf calls whose arguments do not align with the format string.`
}

// typecheck is the check for the go vet command
//
type typecheck struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g typecheck) Name() string { return "typecheck" }

// Weight returns the weight this check has in the overall average
func (g typecheck) Weight() float64 { return .05 }

// Percentage returns the percentage of .go files that pass go vet
func (g typecheck) Percentage() (float64, []model.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=typecheck",
		"--allow-parallel-runners",
	}
	return cmdHelper(g.Dir, g.Filenames, command)
}

// Desc returns the description of go lint
func (g typecheck) Description() string {
	return `<code>go vet</code> examines Go source code and reports suspicious constructs, 
	such as Printf calls whose arguments do not align with the format string.`
}

// unused is the check for the go vet command
//
type unused struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g unused) Name() string { return "unused" }

// Weight returns the weight this check has in the overall average
func (g unused) Weight() float64 { return .05 }

// Percentage returns the percentage of .go files that pass go vet
func (g unused) Percentage() (float64, []model.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=unused",
		"--allow-parallel-runners",
	}
	return cmdHelper(g.Dir, g.Filenames, command)
}

// Desc returns the description of go lint
func (g unused) Description() string {
	return `<code>go vet</code> examines Go source code and reports suspicious constructs, 
	such as Printf calls whose arguments do not align with the format string.`
}

// varcheck is the check for the go vet command
//
type varcheck struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g varcheck) Name() string { return "varcheck" }

// Weight returns the weight this check has in the overall average
func (g varcheck) Weight() float64 { return .05 }

// Percentage returns the percentage of .go files that pass go vet
func (g varcheck) Percentage() (float64, []model.FileSummary, error) {
	command := []string{
		"golangci-lint", "run",
		"--out-format=json",
		"--deadline=180s",
		"--disable-all",
		"--enable=varcheck",
		"--allow-parallel-runners",
	}
	return cmdHelper(g.Dir, g.Filenames, command)
}

// Desc returns the description of go lint
func (g varcheck) Description() string {
	return `<code>go vet</code> examines Go source code and reports suspicious constructs, 
	such as Printf calls whose arguments do not align with the format string.`
}
