package url

import (
	"context"
)

func (r *Repository) Create(ctx context.Context, originalUrl string, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.urls[id] = originalUrl

	return r.save()
}
