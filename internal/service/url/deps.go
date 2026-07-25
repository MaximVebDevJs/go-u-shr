package url

import (
	"context"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
)

type UrlRepository interface {
	Get(ctx context.Context, id string) (string, error)
	Create(ctx context.Context, userID string, originalURL string, id string) (string, error)
	CreateBatch(ctx context.Context, userID string, records []model.BatchRecord) error
	GetUserURLs(ctx context.Context, userID string) ([]model.UserURL, error)
	MarkDeleted(ctx context.Context, userID string, ids []string) error
	Ping(ctx context.Context) error
}
