package repository

import "github.com/pkg/errors"

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
)

// Init with specify DB type
func Init(db DBType) (err error) {
	switch db {
	case Redis:
		_repo, err = NewRedisRepo()
	case Badger:
		_repo, err = NewBadgerRepo("./.badger")
	case Unknown:
		fallthrough
	default:
		_repo, err = NewBadgerRepo("./.badger")
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
