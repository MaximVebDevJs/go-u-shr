package url

import (
	"context"
	"fmt"
)

func (s *service) Get(ctx context.Context, id string) (string, error) {
	url, err := s.urlRepo.Get(ctx, id)
	if err != nil {
		return "", fmt.Errorf("получить utl: %w", err)
	}

	return url, nil
}
