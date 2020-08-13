package vcshelper

import (
	"testing"
)

func Test_builtinToolVCS_Download(t *testing.T) {
	cfgs := []*VCSOption{
		{Host: "github.com", Prefix: "git"},
		{Host: "git.medlinker.com", Prefix: "medgit"},
	}
	vcs := newBuiltinToolVCS(cfgs)

	type args struct {
		remoteURI string
		localDir  string
		branch    string
	}
	tests := []struct {
		name         string
		args         args
		wantRepoRoot string
		wantErr      bool
	}{
		{
			name: "case 0",
			args: args{
				remoteURI: "git.medlinker.com/yeqown/micro-server-template",
				localDir:  "./testdata",
			},
			wantRepoRoot: "testdata/git.medlinker.com/yeqown/micro-server-template",
			wantErr:      false,
		},
		{
			name: "case 1",
			args: args{
				remoteURI: "github.com/yeqown/micro-server-demo",
				localDir:  "./testdata",
			},
			wantRepoRoot: "testdata/github.com/yeqown/micro-server-demo",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepoRoot, err := vcs.Download(tt.args.remoteURI, tt.args.localDir, tt.args.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinToolVCS.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRepoRoot != tt.wantRepoRoot {
				t.Errorf("builtinToolVCS.Download() = %v, want %v", gotRepoRoot, tt.wantRepoRoot)
			}
		})
	}
}
