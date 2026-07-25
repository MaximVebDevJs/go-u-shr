package url

import (
	"context"
	"errors"
	"fmt"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
)

func (s *service) Get(ctx context.Context, id string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("получить url: %w", err)
	}

	url, err := s.urlRepo.Get(ctx, id)
	if err != nil {
		// Маппим инфраструктурные ошибки в сервисные, чтобы HTTP-слой выбрал публичный статус.
		if errors.Is(err, model.ErrURLNotFound) {
			return "", ErrUrlNotFound
		}
		if errors.Is(err, model.ErrURLDeleted) {
			return "", ErrURLDeleted
		}

		return "", fmt.Errorf("получить url: %w", err)
	}

	return url, nil
}
