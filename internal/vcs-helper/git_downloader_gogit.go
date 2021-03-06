package vcshelper

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

//
//import (
//	"fmt"
//	"io/ioutil"
//	"os"
//	"path/filepath"
//	"strings"
//
//	"github.com/yeqown/goreportcard/internal/helper"
//
//	gogit "github.com/go-git/go-git/v5"
//	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
//	"github.com/pkg/errors"
//	"github.com/yeqown/log"
//)
//
//// gitDownloader to clone repo withs git ssh request
//type gitDownloader struct {
//	publicKeys  map[string]*gitssh.PublicKeys // map[host]gitssh.publickeys
//	gitPrefixes map[string]string             // map[host]prefix
//}
//
//// NewGitDownloader with ssh configs
//func NewGitDownloader(cfgs []*VCSOption) IDownloader {
//	gitd := gitDownloader{
//		publicKeys:  make(map[string]*gitssh.PublicKeys),
//		gitPrefixes: make(map[string]string),
//	}
//
//	// load pem file with host
//	for _, v := range cfgs {
//		path, err := filepath.Abs(v.PrivateKeyPath)
//		if err != nil {
//			log.Errorf("NewGitDownloader failed to get Abs path, err=%v", err)
//			continue
//		}
//
//		// load private key
//		pemBytes, err := ioutil.ReadFile(path)
//		if err != nil {
//			log.Errorf("NewGitDownloader failed to open private key file, err=%v", err)
//			continue
//		}
//
//		// Note: username should be PREFIX of ssh clone URL
//		// PREFIX@github.com:OWNER/PROJECT
//		auth, err := gitssh.NewPublicKeys(v.Prefix, []byte(pemBytes), "")
//		if err != nil {
//			log.Errorf("NewGitDownloader failed to NewPublicKeys, err=%v", err)
//			continue
//		}
//
//		gitd.publicKeys[v.Host] = auth
//		gitd.gitPrefixes[v.Host] = v.Prefix
//	}
//
//	return gitd
//}
//
//// TODO: with retry less than 3 times
//// @return repo = github.com/owner/xxx
//// @return error
//func (d gitDownloader) Download(repoURL, localDir, branch string) (string, error) {
//	return d.clone(repoURL, localDir)
//}
//
//// clone to clone repo from remote server
//// it will use ssh public key to clone, config is loaded from config file
//func (d gitDownloader) clone(repoURL string, localDir string) (string, error) {
//	outs, err := hdlRepoURL(repoURL)
//	if err != nil {
//		log.Errorf("could hdl repoURL=%s, err=%v", repoURL, err)
//		return localDir, errors.Wrap(err, "gitDownload.clone failed to hdlRepoURL")
//	}
//	host, owner, repoName := outs[0], outs[1], outs[2]
//
//	// make sure the path has exists
//	localDir = filepath.Join(localDir, host, owner, repoName)
//	if err := helper.EnsurePath(localDir); err != nil {
//		return localDir, errors.Wrap(err, "gitDownload.clone.EnsurePath")
//	}
//
//	// get sshConfig and prefix of git server to clone
//	auth, ok := d.publicKeys[host]
//	if !ok {
//		log.Errorf("gitDownload.clone failed to get sshConfig of host=%s", host)
//		return localDir, errors.New("gitDownload.clone no such host config")
//	}
//	prefix := d.gitPrefixes[host]
//	log.Infof("starting clone with url=%s", wrapRepoURL(prefix, host, owner, repoName))
//
//	// start clone repo
//	repo, err := gogit.PlainClone(localDir, false, &gogit.CloneOptions{
//		URL:        wrapRepoURL(prefix, host, owner, repoName),
//		Auth:       auth,
//		Depth:      1,
//		RemoteName: "origin",
//		Progress:   os.Stdout,
//	})
//	if err != nil {
//		if err == gogit.ErrRepositoryAlreadyExists {
//			// true: repo exist err should not raise
//			return localDir, nil
//		}
//
//		return localDir, errors.Wrap(err, "gitDownloader.clone failed")
//	}
//
//	// TODO: do something with repo
//	_ = repo
//
//	return localDir, nil
//}

// hdlRepoURL does following works:
// @repoURL = github.com/owner/repo
// @return []string{github.com, owner, repo}
// @return error
func hdlRepoURL(repoURL string) ([]string, error) {
	outs := strings.Split(repoURL, "/")
	if len(outs) != 3 {
		return nil, errors.New("hdlRepoURL recv an invalid repoURL")
	}
	return outs, nil
}

// wrapRepoURL to assemble element as "prefix@host:owner/repoName"
func wrapRepoURL(prefix, host, owner, repoName string) string {
	return fmt.Sprintf("%s@%s:%s/%s.git", prefix, host, owner, repoName)
}
