package model

import "time"

// Error contains the line number and the reason for
// an error output from a command
type Error struct {
	LineNumber  int    `json:"line_number"`
	ErrorString string `json:"error_string"`
}

// FileSummary contains the filename, location of the file
// on GitHub, and all of the errors related to the file
type FileSummary struct {
	Filename string  `json:"filename"`
	FileURL  string  `json:"file_url"`
	Errors   []Error `json:"errors"`
}

// AddError adds an Error to FileSummary
func (fs *FileSummary) AddError(err Error) {
	fs.Errors = append(fs.Errors, err)
}

// Score represents the result of a single check
type Score struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	FileSummaries []FileSummary `json:"file_summaries"`
	Weight        float64       `json:"weight"`
	Percentage    float64       `json:"percentage"`
	Error         string        `json:"error"`
}

// LintResult report structure of a lint process to some repository
type LintResult struct {
	Checks               []Score   `json:"checks"`
	Average              float64   `json:"average"`
	Grade                Grade     `json:"grade"`
	Files                int       `json:"files"`
	Issues               int       `json:"issues"`
	Repo                 string    `json:"repo"`
	ResolvedRepo         string    `json:"resolvedRepo"`
	LastRefresh          time.Time `json:"last_refresh"`
	LastRefreshFormatted string    `json:"formatted_last_refresh"`
	LastRefreshHumanized string    `json:"humanized_last_refresh"`
}

// ChecksResult represents the combined result of multiple checks
type ChecksResult struct {
	Scores  []Score `json:"checks"`
	Average float64 `json:"average"`
	Grade   Grade   `json:"GradeFromPercentage"`
	Files   int     `json:"files"`
	Issues  int     `json:"issues"`
}

// ByWeight implements sorting for checks by weight descending
type ByWeight []Score

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].Weight > a[j].Weight }
