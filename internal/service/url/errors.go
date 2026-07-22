package url

import "errors"

var (
	ErrInvalidUrl  = errors.New("невалидный url")
	ErrInvalidJSON = errors.New("невалидный JSON")
	ErrUrlNotFound = errors.New("url не найден")
)
