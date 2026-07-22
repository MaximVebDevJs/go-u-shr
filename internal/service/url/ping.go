package url

import (
	"context"
)

func (s *service) Ping(ctx context.Context) error {
	return s.urlRepo.Ping(ctx)
}
