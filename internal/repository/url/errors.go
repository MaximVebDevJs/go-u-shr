package url

import "errors"

var (
	// ErrNotFound возвращается, когда короткий alias отсутствует в хранилище.
	ErrNotFound = errors.New("не найдено")

	// ErrAlreadyExists возвращается при попытке перезаписать существующий alias.
	ErrAlreadyExists = errors.New("alias уже существует")
)
