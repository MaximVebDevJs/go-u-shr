package url

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const pgUniqueViolationCode = "23505"

func mapCreateError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
		if isUUIDConstraint(pgErr) {
			return ErrAlreadyExists
		}
	}

	return fmt.Errorf("создать url: %w", err)
}

func isUUIDConstraint(pgErr *pgconn.PgError) bool {
	constraint := strings.ToLower(pgErr.ConstraintName)
	column := strings.ToLower(pgErr.ColumnName)

	if strings.Contains(constraint, "uuid") {
		return true
	}

	return column == "uuid"
}

func mapGetError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return fmt.Errorf("получить url: %w", err)
}
