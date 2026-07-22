//go:build integration

package url_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	urlRepo "github.com/MaximVebDevJs/go-u-shr/internal/repository/url"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	schemaOnce sync.Once
	schemaErr  error
)

// ensureSchema выполняет Migrate один раз на процесс: параллельный CREATE TABLE
// гоняется на системном каталоге PostgreSQL (pg_class_relname_nsp_index).
func ensureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	schemaOnce.Do(func() {
		schemaErr = urlRepo.Migrate(ctx, pool)
	})

	return schemaErr
}

func setupTestRepository(t *testing.T) (*urlRepo.Repository, *pgxpool.Pool) {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_DSN")
	}
	require.NotEmpty(t, dsn, "TEST_DATABASE_DSN или DATABASE_DSN обязателен для integration-тестов")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	require.NoError(t, pool.Ping(ctx))
	require.NoError(t, ensureSchema(ctx, pool))

	_, err = pool.Exec(ctx, "TRUNCATE TABLE urls RESTART IDENTITY")
	require.NoError(t, err)

	return urlRepo.New(pool), pool
}

func TestRepositoryPing(t *testing.T) {
	repo, _ := setupTestRepository(t)

	require.NoError(t, repo.Ping(context.Background()))
}

func TestRepositoryCreateAndGet(t *testing.T) {
	ctx := context.Background()
	repo, _ := setupTestRepository(t)

	const (
		id          = "abc12345"
		originalURL = "https://example.com/page"
	)

	require.NoError(t, repo.Create(ctx, originalURL, id))

	got, err := repo.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, originalURL, got)
}

func TestRepositoryGetNotFound(t *testing.T) {
	ctx := context.Background()
	repo, _ := setupTestRepository(t)

	_, err := repo.Get(ctx, "missing")
	require.Error(t, err)
	assert.ErrorIs(t, err, urlRepo.ErrNotFound)
}

func TestRepositoryCreateDuplicateReturnsError(t *testing.T) {
	ctx := context.Background()
	repo, _ := setupTestRepository(t)

	const id = "dup12345"

	require.NoError(t, repo.Create(ctx, "https://example.com/a", id))

	err := repo.Create(ctx, "https://example.com/b", id)
	require.Error(t, err)
	assert.ErrorIs(t, err, urlRepo.ErrAlreadyExists)

	got, getErr := repo.Get(ctx, id)
	require.NoError(t, getErr)
	assert.Equal(t, "https://example.com/a", got)
}

func TestRepositoryCreateConcurrent(t *testing.T) {
	ctx := context.Background()
	repo, _ := setupTestRepository(t)

	const workers = 32

	var wg sync.WaitGroup
	wg.Add(workers)

	errCh := make(chan error, workers)

	for i := range workers {
		go func(n int) {
			defer wg.Done()

			id := fmt.Sprintf("worker%02d", n)
			originalURL := "https://example.com/" + id

			errCh <- repo.Create(ctx, originalURL, id)
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err)
	}
}

func TestMigrateCreatesTable(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_DSN")
	}
	require.NotEmpty(t, dsn, "TEST_DATABASE_DSN или DATABASE_DSN обязателен для integration-тестов")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	require.NoError(t, pool.Ping(ctx))
	require.NoError(t, ensureSchema(ctx, pool))

	var exists bool
	err = pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'urls'
		)`).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists)
}
