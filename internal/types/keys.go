package types

import "strings"

const (
	MasterBranch  = "master"
	RepoFormKey   = "repo"
	BranchFormKey = "branch"
)

type RepoReportParam struct {
	name     string
	branch   string
	identity string
}

func NewRepoParam(repo, branch string) *RepoReportParam {
	return &RepoReportParam{
		name:     repo,
		branch:   branch,
		identity: "",
	}
}

func (p RepoReportParam) RepoIdentity() string {
	if p.identity != "" {
		return p.identity
	}
	p.identity = p.name + "@" + p.branch
	return p.identity
}

func (p RepoReportParam) Branch() string {
	return p.branch
}

func (p RepoReportParam) Repo() string {
	return p.name
}

func ParseRepoIdentity(identity string) (repo, branch string) {
	arr := strings.Split(identity, "@")
	if len(arr) != 2 {
		return identity, ""
	}

	return arr[0], arr[1]
}
