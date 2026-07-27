package url

import "errors"

var (
	ErrInvalidUrl       = errors.New("невалидный url")
	ErrInvalidJSON      = errors.New("невалидный JSON")
	ErrUrlNotFound      = errors.New("url не найден")
	ErrURLDeleted       = errors.New("url удалён")
	ErrURLAlreadyExists = errors.New("url уже существует")

	// ErrBatchTooLarge возвращается, когда в одном запросе прислали больше URL, чем разрешено.
	ErrBatchTooLarge = errors.New("слишком большой batch")

	// ErrDeleteQueueFull возвращается, когда очередь асинхронного удаления переполнена
	// и задачу некуда поставить без блокировки запроса.
	ErrDeleteQueueFull = errors.New("очередь удаления переполнена")
)

// DuplicateURLsError содержит original_url, вызвавшие конфликт в batch.
type DuplicateURLsError struct {
	URLs []string
}

func (e *DuplicateURLsError) Error() string {
	return "дублирующиеся original url"
}
