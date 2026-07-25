package model

import "errors"

var (
	// ErrURLNotFound возвращается, когда короткий alias отсутствует в хранилище.
	ErrURLNotFound = errors.New("не найдено")

	// ErrAliasAlreadyExists возвращается при попытке перезаписать существующий alias.
	ErrAliasAlreadyExists = errors.New("alias уже существует")

	// ErrOriginalURLExists возвращается, когда original_url уже есть в хранилище.
	ErrOriginalURLExists = errors.New("original url уже существует")

	// ErrURLDeleted возвращается, когда alias найден, но помечен удалённым.
	ErrURLDeleted = errors.New("url удалён")
)

// BatchRecord описывает одну доменную запись для пакетного сохранения URL.
type BatchRecord struct {
	OriginalURL string
	ID          string
}

// UserURL описывает URL, созданный конкретным пользователем.
type UserURL struct {
	ShortID     string
	OriginalURL string
}

// DuplicateOriginalURLsError содержит список original_url, вызвавших конфликт.
type DuplicateOriginalURLsError struct {
	URLs []string
}

// Error возвращает текстовое описание конфликта original_url.
func (e *DuplicateOriginalURLsError) Error() string {
	return "дублирующиеся original url"
}
