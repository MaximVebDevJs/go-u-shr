package url

import (
	"sync"
)

type Repository struct {
	mu       sync.RWMutex
	urls     map[string]string
	filePath string
}

func New(filePath string) *Repository {
	r := &Repository{
		urls:     make(map[string]string),
		filePath: filePath,
	}

	_ = r.load()

	return r
}
