package url

import (
	"context"
)

func (r *Repository) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}
