package url

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	// deleteChunkSize ограничивает длину массива id в одном batch update.
	deleteChunkSize     = 100
	deleteQueueSize     = 1024
	deleteFlushInterval = 100 * time.Millisecond
	deleteFlushTimeout  = 5 * time.Second
	deleteFlushAttempts = 3
	deleteRetryDelay    = 200 * time.Millisecond
)

// errServiceStopped возвращается, когда задачу удаления некому обработать: worker уже остановлен.
var errServiceStopped = errors.New("service остановлен")

type deleteTask struct {
	userID string
	ids    []string
}

// DeleteUrls принимает URL пользователя на асинхронное удаление и не ждёт выполнения SQL update.
func (s *service) DeleteUrls(ctx context.Context, userID string, ids []string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("удалить urls: %w", err)
	}

	// Размер запроса ограничен на транспорте через MaxBytesReader, а батчи режет worker,
	// поэтому здесь проверяем только бизнес-инвариант: есть ли что удалять.
	uniqueIDs := normalizeDeleteIDs(ids)
	if len(uniqueIDs) == 0 {
		return fmt.Errorf("удалить urls: %w", ErrInvalidUrl)
	}

	task := deleteTask{userID: userID, ids: uniqueIDs}

	select {
	case <-ctx.Done():
		return fmt.Errorf("поставить urls на удаление: %w", ctx.Err())
	case <-s.stop:
		return fmt.Errorf("поставить urls на удаление: %w", errServiceStopped)
	case s.deleteQueue <- task:
		return nil
	}
}

func normalizeDeleteIDs(ids []string) []string {
	// Дедупликация делает DELETE идемпотентным и не раздувает batch update при повторяющихся id.
	seen := make(map[string]struct{}, len(ids))
	result := make([]string, 0, len(ids))

	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}

		if _, ok := seen[id]; ok {
			continue
		}

		seen[id] = struct{}{}
		result = append(result, id)
	}

	return result
}

func (s *service) runDeleteWorker() {
	defer s.wg.Done()

	ticker := time.NewTicker(deleteFlushInterval)
	defer ticker.Stop()

	pending := make(map[string]map[string]struct{})

	for {
		select {
		case <-s.stop:
			drainDeleteQueue(s.deleteQueue, pending)
			s.flushDeleteBatch(pending)
			return
		case task := <-s.deleteQueue:
			addDeleteTask(pending, task)
			if deleteBatchSize(pending) >= deleteChunkSize {
				s.flushDeleteBatch(pending)
				pending = make(map[string]map[string]struct{})
			}
		case <-ticker.C:
			if deleteBatchSize(pending) == 0 {
				continue
			}

			s.flushDeleteBatch(pending)
			pending = make(map[string]map[string]struct{})
		}
	}
}

func drainDeleteQueue(queue <-chan deleteTask, batch map[string]map[string]struct{}) {
	// На shutdown HTTP уже остановлен, поэтому можно добрать всё, что успели принять до Close.
	for {
		select {
		case task := <-queue:
			addDeleteTask(batch, task)
		default:
			return
		}
	}
}

func addDeleteTask(batch map[string]map[string]struct{}, task deleteTask) {
	ids, ok := batch[task.userID]
	if !ok {
		ids = make(map[string]struct{}, len(task.ids))
		batch[task.userID] = ids
	}

	for _, id := range task.ids {
		ids[id] = struct{}{}
	}
}

func deleteBatchSize(batch map[string]map[string]struct{}) int {
	total := 0
	for _, ids := range batch {
		total += len(ids)
	}

	return total
}

func (s *service) flushDeleteBatch(batch map[string]map[string]struct{}) {
	for userID, idsSet := range batch {
		ids := idsFromSet(idsSet)
		for len(ids) > 0 {
			chunkSize := min(len(ids), deleteChunkSize)
			chunk := ids[:chunkSize]
			ids = ids[chunkSize:]

			s.markDeleted(userID, chunk)
		}
	}
}

// markDeleted повторяет batch update при ошибке: клиент уже получил 202 и не узнает о сбое,
// а сам update идемпотентен, поэтому повтор безопасен.
func (s *service) markDeleted(userID string, ids []string) {
	for attempt := 1; attempt <= deleteFlushAttempts; attempt++ {
		// Flush использует собственный timeout: request context уже мог завершиться после ответа 202.
		ctx, cancel := context.WithTimeout(context.Background(), deleteFlushTimeout)
		err := s.urlRepo.MarkDeleted(ctx, userID, ids)
		cancel()

		if err == nil {
			s.logger.Info("ссылки помечены удалёнными",
				zap.String("user_id", userID),
				zap.Int("batch_size", len(ids)),
				zap.Int("attempt", attempt),
			)

			return
		}

		s.logger.Error("ошибка асинхронного удаления ссылок",
			zap.Error(err),
			zap.String("user_id", userID),
			zap.Int("batch_size", len(ids)),
			zap.Int("attempt", attempt),
		)

		if attempt == deleteFlushAttempts {
			return
		}

		// На shutdown повторы не ждём: время остановки ограничено shutdownTimeout.
		select {
		case <-s.stop:
			return
		case <-time.After(deleteRetryDelay):
		}
	}
}

func idsFromSet(idsSet map[string]struct{}) []string {
	ids := make([]string, 0, len(idsSet))
	for id := range idsSet {
		ids = append(ids, id)
	}

	return ids
}
