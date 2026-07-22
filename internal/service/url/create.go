package url

import (
	"context"
	"errors"
	"fmt"

	urlRepo "github.com/MaximVebDevJs/go-u-shr/internal/repository/url"
)

// maxIDGenerationAttempts — сколько раз пробуем сгенерировать ID при коллизии.
const maxIDGenerationAttempts = 5

func (s *service) Create(ctx context.Context, originalURL string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("создать url: %w", err)
	}

	// Валидация до записи в хранилище — защита от open redirect через Location.
	if err := validateURL(originalURL); err != nil {
		return "", err
	}

	for range maxIDGenerationAttempts {
		id := generateShortID()

		err := s.urlRepo.Create(ctx, originalURL, id)
		if err == nil {
			return s.baseURL + "/" + id, nil
		}

		// Коллизия ID — пробуем другой alias.
		if errors.Is(err, urlRepo.ErrAlreadyExists) {
			continue
		}

		return "", fmt.Errorf("сохранение url: %w", err)
	}

	return "", fmt.Errorf("создать url: не удалось сгенерировать уникальный id за %d попыток", maxIDGenerationAttempts)
}
