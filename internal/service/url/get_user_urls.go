package url

import (
	"context"
	"fmt"
)

// GetUserURLs возвращает все сокращенные URL пользователя.
func (s *Service) GetUserURLs(ctx context.Context, userID string) ([]UserURL, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("получить ссылки пользователя: %w", err)
	}

	urls, err := s.urlRepo.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("получить ссылки пользователя: %w", err)
	}

	result := make([]UserURL, 0, len(urls))
	for _, item := range urls {
		result = append(result, UserURL{
			ShortURL:    s.baseURL + "/" + item.ShortID,
			OriginalURL: item.OriginalURL,
		})
	}

	return result, nil
}
