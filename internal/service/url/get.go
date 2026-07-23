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
		// Маппим инфраструктурную ошибку в доменную, чтобы HTTP-слой отдал 404, а не 500.
		if errors.Is(err, model.ErrURLNotFound) {
			return "", ErrUrlNotFound
		}

		return "", fmt.Errorf("получить url: %w", err)
	}

	return url, nil
}
