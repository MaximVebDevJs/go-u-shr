package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func main() {
	cfg := config.Load()

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		slog.Error("не удалось инициализировать логер", "error", err)
		os.Exit(1)
	}
	defer log.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startupCtx, cancelStartup := context.WithTimeout(ctx, startupTimeout)
	defer cancelStartup()

	pool, err := pgxpool.New(startupCtx, cfg.DatabaseDSN)
	if err != nil {
		slog.Error("создание пула соединений", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err = pool.Ping(startupCtx); err != nil {
		slog.Error("проверка удалось ли подключиться к базе данных", "error", err)
		os.Exit(1)
	}

	if err = urlRepo.Migrate(startupCtx, pool); err != nil {
		slog.Error("миграция базы данных", "error", err)
		os.Exit(1)
	}
	cancelStartup()

	if err = run(ctx, cfg, log, pool); err != nil {
		slog.Error("ошибка запуска urlShortener сервиса", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg *config.Config, log *zap.Logger, pool *pgxpool.Pool) error {
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
