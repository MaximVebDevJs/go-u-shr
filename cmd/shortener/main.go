package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/MaximVebDevJs/go-u-shr/internal/auth"
	"github.com/MaximVebDevJs/go-u-shr/internal/config"
	"github.com/MaximVebDevJs/go-u-shr/internal/logger"
	urlRepo "github.com/MaximVebDevJs/go-u-shr/internal/repository/url"
	"github.com/MaximVebDevJs/go-u-shr/pkg/app"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	startupTimeout  = 10 * time.Second
	shutdownTimeout = 10 * time.Second
	readTimeout     = 15 * time.Second
	writeTimeout    = 15 * time.Second
	idleTimeout     = 60 * time.Second
)

const exitFailure = 1

func main() {
	// os.Exit не выполняет defer, поэтому вся инициализация живёт в run.
	os.Exit(run())
}

func run() int {
	cfg := config.Load()

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "инициализация логера: %v\n", err)

		return exitFailure
	}
	defer func() { _ = log.Sync() }()

	if err = ensureAuthSecret(cfg, log); err != nil {
		log.Error("подготовка секрета авторизации", zap.Error(err))

		return exitFailure
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := newPool(ctx, cfg, log)
	if err != nil {
		log.Error("инициализация базы данных", zap.Error(err))

		return exitFailure
	}
	defer pool.Close()

	if err = serve(ctx, cfg, log, pool); err != nil {
		log.Error("работа urlShortener сервиса", zap.Error(err))

		return exitFailure
	}

	return 0
}

// ensureAuthSecret подставляет одноразовый секрет, если он не задан конфигурацией.
func ensureAuthSecret(cfg *config.Config, log *zap.Logger) error {
	if strings.TrimSpace(cfg.AuthSecret) != "" {
		return nil
	}

	secret, err := auth.GenerateSecret()
	if err != nil {
		return fmt.Errorf("сгенерировать секрет авторизации: %w", err)
	}

	cfg.AuthSecret = secret

	log.Warn("AUTH_SECRET не задан: сгенерирован временный секрет, " +
		"выданные cookie не переживут перезапуск процесса")

	return nil
}

func newPool(ctx context.Context, cfg *config.Config, log *zap.Logger) (*pgxpool.Pool, error) {
	startupCtx, cancel := context.WithTimeout(ctx, startupTimeout)
	defer cancel()

	pool, err := pgxpool.New(startupCtx, cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("создание пула соединений: %w", err)
	}

	if err = pool.Ping(startupCtx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("проверка соединения с базой данных: %w", err)
	}

	if err = urlRepo.Migrate(startupCtx, pool); err != nil {
		pool.Close()

		return nil, fmt.Errorf("миграция базы данных: %w", err)
	}

	log.Info("подключение к базе данных установлено")

	return pool, nil
}

func serve(ctx context.Context, cfg *config.Config, log *zap.Logger, pool *pgxpool.Pool) error {
	// Создать HTTP-обработчик через фабрику приложения.
	handler, cleanup, err := app.NewHTTPHandler(cfg, log, pool)
	if err != nil {
		return fmt.Errorf("создание HTTP handler: %w", err)
	}
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Service cleanup останавливает фоновые worker-ы после того, как HTTP перестал принимать новые запросы.
		if cleanupErr := cleanup(cleanupCtx); cleanupErr != nil {
			log.Error("ошибка остановки фоновых задач", zap.Error(cleanupErr))
		}
	}()

	httpServer := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	serverErr := make(chan error, 1)

	// Запускаем HTTP-сервер в отдельной горутине.
	go func() {
		log.Info("HTTP сервер запущен", zap.String("address", cfg.ServerAddr))

		if listenErr := httpServer.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("HTTP сервер: %w", listenErr)
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
	}

	log.Info("остановка urlShortener")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if shutdownErr := httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
		return fmt.Errorf("остановка HTTP сервера: %w", shutdownErr)
	}

	log.Info("urlShortener остановлен")

	return nil
}
