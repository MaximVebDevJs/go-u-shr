package url

import (
	"context"
)

func (r *Repository) Get(ctx context.Context, id string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rUrl, ok := r.urls[id]
	if !ok {
		return "", ErrNotFound
	}

	return rUrl, nil
}
