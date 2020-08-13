package vcshelper

import "fmt"

var (
	_downloader IDownloader
)

type IDownloader interface {
	// Download to copy git repo to local folders
	// @repoRoot is relative path to repo
	// @err error
	Download(remoteURL, localDir, branch string) (repoRoot string, err error)
}

// TODO: 并发安全
// GetDownloader get the builtin git downloader variable
func GetDownloader() IDownloader {
	return _downloader
}

type VCSType uint8

const (
	Unknown VCSType = iota
	BuiltinTool
	GoGit
)

// VCSOption ssh clone public key config
type VCSOption struct {
	Host           string // host of git server
	PrivateKeyPath string // private key pem path
	Prefix         string // prefix of git server. refer prefix@host:owner/repoName
}

// Init downloader
// provide an option to switch the initial downloader between go-VCS and git
func Init(vcs VCSType, opts []*VCSOption) error {
	switch vcs {
	case Unknown:
		fallthrough
	default:
		fallthrough
	//case GoGit:
	// _downloader = NewGitDownloader(opts)
	//fallthrough
	case BuiltinTool:
		_downloader = newBuiltinToolVCS(opts)
	}

	if _downloader != nil {
		return nil
	}

	return fmt.Errorf("invalid vcs type: %d", vcs)
}
