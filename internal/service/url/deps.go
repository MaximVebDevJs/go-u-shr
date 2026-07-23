package url

import (
	"context"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
)

type UrlRepository interface {
	Get(ctx context.Context, id string) (string, error)
	Create(ctx context.Context, originalURL string, id string) (string, error)
	CreateBatch(ctx context.Context, records []model.BatchRecord) error
	Ping(ctx context.Context) error
}
