package url

import (
	"context"

	errs "github.com/MaximVebDevJs/go-u-shr/internal/errors"
	"github.com/MaximVebDevJs/go-u-shr/internal/model"
)

func (r *Repository) Get(_ context.Context, url string) (model.Url, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rUrl, ok := r.urls[url]
	if !ok {
		return model.Url{}, errs.ErrUrlNotFound
	}

	return rUrl, nil
}
