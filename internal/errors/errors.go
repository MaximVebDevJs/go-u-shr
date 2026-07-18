package errs

import "errors"

var (
	ErrUrlNotFound = errors.New("url не найден")
	ErrInvalidUrl  = errors.New("невалидный url")
	ErrInvalidJSON = errors.New("невалидный JSON")
	ErrEmptyId     = errors.New("пустой id")
)
