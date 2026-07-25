package url

import "errors"

var (
	ErrInvalidUrl       = errors.New("невалидный url")
	ErrInvalidJSON      = errors.New("невалидный JSON")
	ErrUrlNotFound      = errors.New("url не найден")
	ErrURLDeleted       = errors.New("url удалён")
	ErrURLAlreadyExists = errors.New("url уже существует")
)

// DuplicateURLsError содержит original_url, вызвавшие конфликт в batch.
type DuplicateURLsError struct {
	URLs []string
}

func (e *DuplicateURLsError) Error() string {
	return "дублирующиеся original url"
}
