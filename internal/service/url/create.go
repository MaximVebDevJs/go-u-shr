package url

import (
	"context"
	"fmt"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
)

func (s *service) Create(ctx context.Context, req model.CreateUrlRequest) (model.Url, error) {
	// логика
	url := model.Url{
		Url: req.Url,
	}

	if _, err := s.urlRepo.Create(ctx, url); err != nil {
		return model.Url{}, fmt.Errorf("сохранение url: %w", err)
	}

	return url, nil
}
