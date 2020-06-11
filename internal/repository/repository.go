package repository

import "github.com/pkg/errors"

// IRepository KV DB operator
type IRepository interface {
	// Get value with key
	Get(key []byte) ([]byte, error)

	// Set value with key
	Update(key, value []byte) error

	// Close DB connection
	Close()
}

var (
	_repo IRepository

	// ErrKeyNotFound .
	ErrKeyNotFound = errors.New("key not found")
)

// Init .
func Init() (err error) {
	_repo, err = NewBadgerRepo("./.badger")
	//_repo, err = NewRedisRepo()
	if err != nil {
		return err
	}

	return nil
}

// GetRepo .
func GetRepo() IRepository {
	return _repo
}
