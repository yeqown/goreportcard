package linter

import (
	"reflect"
	"testing"

	"github.com/yeqown/goreportcard/internal/types"
)

// func TestGoFiles(t *testing.T) {
// 	files, skipped, err := visitGoFiles("testfiles/")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	want := []string{"testfiles/a.go", "testfiles/b.go", "testfiles/c.go"}
// 	if !reflect.DeepEqual(files, want) {
// 		t.Errorf("visitGoFiles(%q) = %v, want %v", "testfiles/", files, want)
// 	}

// 	wantSkipped := []string{"testfiles/a.pb.go", "testfiles/a.pb.gw.go"}
// 	if !reflect.DeepEqual(skipped, wantSkipped) {
// 		t.Errorf("visitGoFiles(%q) skipped = %v, want %v", "testfiles/", skipped, wantSkipped)
// 	}
// }

var goToolTests = []struct {
	name      string
	dir       string
	filenames []string
	tool      []string
	percent   float64
	failed    []types.FileSummary
	wantErr   bool
}{
	{"go vet", "testfiles/", []string{"testfiles/a.go", "testfiles/b.go", "testfiles/c.go"}, []string{"go", "tool", "vet"}, 1, []types.FileSummary{}, false},
}

func TestGoTool(t *testing.T) {
	for _, tt := range goToolTests {
		f, fs, err := cmdHelper(tt.dir, tt.filenames, tt.tool)
		if err != nil && !tt.wantErr {
			t.Fatal(err)
		}
		if f != tt.percent {
			t.Errorf("[%s] cmdHelper percent = %f, want %f", tt.name, f, tt.percent)
		}
		if !reflect.DeepEqual(fs, tt.failed) {
			t.Errorf("[%s] cmdHelper failed = %v, want %v", tt.name, fs, tt.failed)
		}
	}
}
