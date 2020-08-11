package repository

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// IRepository .
type IRepository interface {
	Get(key []byte) ([]byte, error)

	Update(key, value []byte) error

	Close()
}

// DBType .
type DBType uint8

const (
	Unknown DBType = iota // unknown
	Redis                 // redis
	Badger                // badger
)

var (
	// _repo instance of IRepository
	_repo IRepository

	// ErrKeyNotFound .
	ErrKeyNotFound = errors.New("key not found")

	// default path to save bader
	_defaultBadgerDBPath = ".badger"
)

func init() {
	home, _ := os.UserHomeDir()
	_defaultBadgerDBPath = filepath.Join(home, _defaultBadgerDBPath)
}

// Init with specify DB type
func Init(db DBType) (err error) {
	switch db {
	case Redis:
		_repo, err = NewRedisRepo()
	case Badger:
		_repo, err = NewBadgerRepo(_defaultBadgerDBPath)
	case Unknown:
		fallthrough
	default:
		_repo, err = NewBadgerRepo(_defaultBadgerDBPath)
	}

	if err != nil {
		return err
	}

	return nil
}

// GetRepo .
func GetRepo() IRepository {
	return _repo
}
