package url

import (
	"context"
	"fmt"

	"github.com/MaximVebDevJs/go-u-shr/internal/model"
)

func (s *service) Get(ctx context.Context, url string) (model.Url, error) {
	rUrl, err := s.urlRepo.Get(ctx, url)
	if err != nil {
		return model.Url{}, fmt.Errorf("получить utl: %w", err)
	}

	return rUrl, nil
}
