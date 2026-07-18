package url

import (
	"context"
	"fmt"
)

func (s *service) Create(ctx context.Context, originalURL string) (string, error) {
	id := generateShortID()

	if err := s.urlRepo.Create(ctx, originalURL, id); err != nil {
		return "", fmt.Errorf("сохранение url: %w", err)
	}

	return s.baseURL + "/" + id, nil
}
