package url

import (
	"sync"
)

type Repository struct {
	mu       sync.RWMutex
	urls     map[string]string
	filePath string
}

func New(filePath string) (*Repository, error) {
	r := &Repository{
		urls:     make(map[string]string),
		filePath: filePath,
	}

	if err := r.load(); err != nil {
		return nil, err
	}

	return r, nil
}
