package url

import "github.com/MaximVebDevJs/go-u-shr/internal/model"

var (
	// ErrNotFound возвращается, когда короткий alias отсутствует в хранилище.
	ErrNotFound = model.ErrURLNotFound

	// ErrAlreadyExists возвращается при попытке перезаписать существующий shortUrl.
	ErrAlreadyExists = model.ErrAliasAlreadyExists

	// ErrOriginalURLExists возвращается, когда original_url уже есть в хранилище.
	ErrOriginalURLExists = model.ErrOriginalURLExists

	// ErrDeleted возвращается, когда запись найдена, но помечена удалённой.
	ErrDeleted = model.ErrURLDeleted
)

// DuplicateOriginalURLsError содержит список original_url, вызвавших конфликт в batch.
type DuplicateOriginalURLsError = model.DuplicateOriginalURLsError
