package url

import (
	"context"
)

func (s *Service) Ping(ctx context.Context) error {
	return s.urlRepo.Ping(ctx)
}
