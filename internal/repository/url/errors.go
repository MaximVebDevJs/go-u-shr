package url

import "errors"

var (
	// ErrNotFound возвращается, когда короткий alias отсутствует в хранилище.
	ErrNotFound = errors.New("не найдено")

	// ErrAlreadyExists возвращается при попытке перезаписать существующий shortUrl.
	ErrAlreadyExists = errors.New("alias уже существует")

	// ErrOriginalURLExists возвращается, когда original_url уже есть в хранилище.
	ErrOriginalURLExists = errors.New("original url уже существует")
)

// DuplicateOriginalURLsError содержит список original_url, вызвавших конфликт в batch.
type DuplicateOriginalURLsError struct {
	URLs []string
}

func (e *DuplicateOriginalURLsError) Error() string {
	return "дублирующиеся original url"
}
