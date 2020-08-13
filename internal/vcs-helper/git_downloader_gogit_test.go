package vcshelper

//
//import (
//	"testing"
//)
//
//func Test_gitDownloader_gitlab(t *testing.T) {
//	cfgs := []*VCSOption{
//		{
//			Host:           "git.medlinker.com",
//			PrivateKeyPath: "/Users/med/.ssh/id_rsa",
//			Prefix:         "medgit",
//		},
//	}
//
//	// medgit@medlinker.com:yeqown/micro-server-template.git
//	gitd := NewGitDownloader(cfgs)
//	root, err := gitd.Download(
//		"git.medlinker.com/yeqown/micro-server-template", "./testdata", "master")
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//	t.Logf("got root=%s", root)
//}
//
//func Test_gitDownloader_github(t *testing.T) {
//	cfgs := []*VCSOption{
//		{
//			Host:           "github.com",
//			PrivateKeyPath: "/Users/med/.ssh/id_rsa",
//			Prefix:         "git",
//		},
//	}
//
//	// git@github.com:yeqown/websocket.git
//	gitd := NewGitDownloader(cfgs)
//	root, err := gitd.Download(
//		"github.com/yeqown/websocket", "./testdata", "master")
//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//
//	t.Logf("got root=%s", root)
//}
