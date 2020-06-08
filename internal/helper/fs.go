package helper

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/yeqown/log"
)

var (
	errPathNotExists = errors.New("path not exists")
)

// exists returns whether the given file or directory exists or not
// from http://stackoverflow.com/a/10510783
func Exists(path string) (ok bool, err error) {
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, errPathNotExists
		}

		return false, errors.Wrap(err, "fs.Exists failed to os.Stat")
	}

	return true, err
}

// EnsurePath make sure the path has been exists.
// it will create if path not exists
func EnsurePath(path string) (err error) {
	var ok bool
	if ok, err = Exists(path); ok {
		return nil
	}

	if err == errPathNotExists {
		// true: not exists then make dirs
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return errors.Wrap(err, "fs.EnsurePath failed to mkdir")
		}

		return nil
	}

	return errors.Wrap(err, "fs.EnsurePath failed to check path")
}

func IsEmptyDir(path string) (empty bool) {
	cnt, err := countFiles(path, true)
	if err != nil {
		// FIXME: TRUE: any error means dir not empty ?
		return false
	}

	return cnt == 0
}

// coutnFiles under path
func countFiles(path string, recursive bool) (cnt int, err error) {
	dirs := make([]string, 1, 10)
	dirs[0] = path

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "countFiles.walkFn")
		}

		if info.IsDir() && recursive {
			dirs = append(dirs, path)
			return nil
		}

		cnt++

		// TODO: maybe need to skip somefiles
		return nil
	}

	for _, dir := range dirs {
		if err = filepath.Walk(dir, walkFn); err != nil {
			log.Warnf("countFiles got an error, dir=%s, err=%v", dir, err)
			err = errors.Wrap(err, "countFiles.Walk")
			return
		}
	}

	return
}
