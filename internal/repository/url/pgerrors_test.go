package url

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapGetError(t *testing.T) {
	t.Parallel()

	t.Run("no rows", func(t *testing.T) {
		t.Parallel()
		assert.ErrorIs(t, mapGetError(pgx.ErrNoRows), ErrNotFound)
	})

	t.Run("other error", func(t *testing.T) {
		t.Parallel()

		err := mapGetError(errors.New("connection reset"))
		require.Error(t, err)
		assert.False(t, errors.Is(err, ErrNotFound))
	})
}

func TestMapCreateError(t *testing.T) {
	t.Parallel()

	t.Run("unique violation", func(t *testing.T) {
		t.Parallel()

		pgErr := &pgconn.PgError{Code: pgUniqueViolationCode}
		assert.ErrorIs(t, mapCreateError(pgErr), ErrAlreadyExists)
	})

	t.Run("other error", func(t *testing.T) {
		t.Parallel()

		err := mapCreateError(errors.New("disk full"))
		require.Error(t, err)
		assert.False(t, errors.Is(err, ErrAlreadyExists))
	})
}
