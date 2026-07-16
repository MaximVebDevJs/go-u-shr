package url

import (
	"context"
)

func (r *Repository) Create(ctx context.Context, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.urls[url] = url

	return nil
}
