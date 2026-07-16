package url

import (
	"context"

	errs "github.com/MaximVebDevJs/go-u-shr/internal/errors"
)

func (r *Repository) Get(ctx context.Context, id string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rUrl, ok := r.urls[id]
	if !ok {
		return "", errs.ErrUrlNotFound
	}

	return rUrl, nil
}
