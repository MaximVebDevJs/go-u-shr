package url

import (
	"context"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
)

func (r *Repository) Create(_ context.Context, url model.Url) (model.Url, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.urls[url.Url] = url

	return url, nil
}
