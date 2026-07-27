package url

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// Service реализует usecase-логику сокращения ссылок и владеет worker-ом асинхронного удаления.
type Service struct {
	urlRepo UrlRepository
	baseURL string
	logger  *zap.Logger

	deleteQueue chan deleteTask
	// stop закрывается в Close и служит сигналом остановки для фоновых worker-ов.
	stop     chan struct{}
	stopOnce sync.Once
	wg       sync.WaitGroup
}

// New создаёт URL service и запускает worker асинхронного удаления.
func New(urlRepo UrlRepository, baseURL string, logger *zap.Logger) *Service {
	if logger == nil {
		logger = zap.NewNop()
	}

	s := &Service{
		urlRepo:     urlRepo,
		baseURL:     baseURL,
		logger:      logger,
		deleteQueue: make(chan deleteTask, deleteQueueSize),
		stop:        make(chan struct{}),
	}

	// Worker принадлежит service: он живёт дольше HTTP-запроса и гарантированно останавливается через Close.
	s.wg.Add(1)
	go s.runDeleteWorker()

	return s
}

// Close останавливает фоновые задачи service и ждёт их завершения до отмены ctx.
func (s *Service) Close(ctx context.Context) error {
	// Once защищает от повторного close(s.stop), если Close вызовут дважды.
	s.stopOnce.Do(func() {
		close(s.stop)
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		s.wg.Wait()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("остановить url service: %w", ctx.Err())
	case <-done:
		return nil
	}
}
