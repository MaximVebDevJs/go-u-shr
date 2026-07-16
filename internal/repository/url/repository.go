package url

import (
	"sync"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
)

type Repository struct {
	mu   sync.RWMutex
	urls map[string]model.Url
}

func New() *Repository {
	return &Repository{
		urls: make(map[string]model.Url),
	}
}
