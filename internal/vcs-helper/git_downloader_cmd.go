package vcshelper

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/yeqown/goreportcard/internal/helper"
	"github.com/yeqown/log"
)

var _ IDownloader = builtinToolVCS{}

type builtinToolVCS struct {
	Cmd         string // git command
	Dir         string // repo dir
	fetchCmd    string // fetch command
	checkoutCmd string // checkout command
	cloneCmd    string // clone command
	pullCmd     string // pull command

	gitPrefixes map[string]string // map[host]prefix
}

// newBuiltinToolVCS .
func newBuiltinToolVCS(cfgs []*VCSOption) IDownloader {
	downloader := builtinToolVCS{
		Cmd:         "git",
		Dir:         "",
		fetchCmd:    "fetch {arg}",
		checkoutCmd: "checkout {branch}",
		cloneCmd:    "clone {remote}",
		pullCmd:     "pull origin {branch}",

		gitPrefixes: make(map[string]string, 4),
	}

	for _, v := range cfgs {
		downloader.gitPrefixes[v.Host] = v.Prefix
	}

	return downloader
}

// 参考 golang.org/x/tools/go/vcs 设计
func (c builtinToolVCS) run(dir string, cmd string, keyval ...string) error {
	_, err := c.run1(dir, cmd, keyval, true)
	return err
}

func (c builtinToolVCS) run1(dir string, cmdline string, keyval []string, verbose bool) ([]byte, error) {
	m := make(map[string]string)
	for i := 0; i < len(keyval); i += 2 {
		m[keyval[i]] = keyval[i+1]
	}
	args := strings.Fields(cmdline)
	for i, arg := range args {
		args[i] = expand(m, arg)
	}

	_, err := exec.LookPath(c.Cmd)
	if err != nil {
		log.Errorf("go: missing %s command.", c.Cmd)
		return nil, err
	}

	cmd := exec.Command(c.Cmd, args...)
	cmd.Dir = dir
	cmd.Env = envForDir(cmd.Dir)

	log.Debugf("cd %s\n", dir)
	log.Debugf("%s %s\n", c.Cmd, strings.Join(args, " "))

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err = cmd.Run()
	out := buf.Bytes()
	if err != nil {
		if verbose {
			log.Errorf("# cd %s; %s %s", dir, c.Cmd, strings.Join(args, " "))
			log.Errorf("%s", out)
		}
		return nil, err
	}
	return out, nil
}

func (c builtinToolVCS) Download(repoURL, localDir, branch string) (string, error) {
	outs, err := hdlRepoURL(repoURL)
	if err != nil {
		log.Errorf("could hdl repoURL=%s, err=%v", repoURL, err)
		return localDir, errors.Wrap(err, "gitDownload.clone failed to hdlRepoURL")
	}
	host, owner, repoName := outs[0], outs[1], outs[2]

	// make sure the path has exists
	localDir = filepath.Join(localDir, host, owner, repoName)
	if err := helper.EnsurePath(localDir); err != nil {
		return localDir, errors.Wrap(err, "gitDownload.clone.EnsurePath")
	}

	// get git prefix
	var (
		prefix string
		ok     bool
	)
	if prefix, ok = c.gitPrefixes[host]; !ok {
		return "", errors.New("gitDownload.clone no such host config")
	}
	remoteURI := wrapRepoURL(prefix, host, owner, repoName)

	if c.shouldCreate(localDir) {
		// TODO: create local folder first
		err = c.run(c.cloneCmd, "remote", remoteURI)
		if err != nil {
			return localDir, err
		}
	} else {
		// fetch all
		err = c.run(c.fetchCmd, "arg", "--all")
		if err != nil {
			return localDir, err
		}

		// TODO: checkout and pull

		// pull
		err = c.run(c.pullCmd, "branch", branch)
	}

	// checkout
	err = c.run(c.checkoutCmd, "branch", branch)
	return localDir, err
}

// expand rewrites s to replace {k} with match[k] for each key k in match.
func expand(match map[string]string, s string) string {
	for k, v := range match {
		s = strings.Replace(s, "{"+k+"}", v, -1)
	}
	return s
}

// envForDir returns a copy of the environment
// suitable for running in the given directory.
// The environment is the current process's environment
// but with an updated $PWD, so that an os.Getwd in the
// child will be faster.
func envForDir(dir string) []string {
	env := os.Environ()
	// Internally we only use rooted paths, so dir is rooted.
	// Even if dir is not rooted, no harm done.
	return mergeEnvLists([]string{"PWD=" + dir}, env)
}

// mergeEnvLists merges the two environment lists such that
// variables with the same name in "in" replace those in "out".
func mergeEnvLists(in, out []string) []string {
NextVar:
	for _, inkv := range in {
		k := strings.SplitAfterN(inkv, "=", 2)[0]
		for i, outkv := range out {
			if strings.HasPrefix(outkv, k) {
				out[i] = inkv
				continue NextVar
			}
		}
		out = append(out, inkv)
	}
	return out
}

// shouldCreate if dir is empty means should Create else Download
func (vcs builtinToolVCS) shouldCreate(localDir string) bool {
	return helper.IsEmptyDir(localDir)
}
