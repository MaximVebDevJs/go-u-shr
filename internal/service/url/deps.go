package url

import (
	"context"

	urlRepo "github.com/MaximVebDevJs/go-u-shr/internal/repository/url"
)

type UrlRepository interface {
	Get(ctx context.Context, id string) (string, error)
	Create(ctx context.Context, originalURL string, id string) (string, error)
	CreateBatch(ctx context.Context, records []urlRepo.BatchRecord) error
	Ping(ctx context.Context) error
}
