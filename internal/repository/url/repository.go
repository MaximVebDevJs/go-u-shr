package url

import (
	"sync"
)

type Repository struct {
	mu   sync.RWMutex
	urls map[string]string
}

func New() *Repository {
	return &Repository{
		urls: make(map[string]string),
	}
}
