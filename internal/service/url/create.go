package url

import (
	"context"
	"fmt"
)

func (s *service) Create(ctx context.Context, originalURL string) (string, error) {
	// логика originalURL

	urlAfterLogic := originalURL

	if err := s.urlRepo.Create(ctx, urlAfterLogic); err != nil {
		return "", fmt.Errorf("сохранение url: %w", err)
	}

	return urlAfterLogic, nil
}
