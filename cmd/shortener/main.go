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

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		slog.Error("создание пула соединений", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		slog.Error("проверка удалось ли подключиться к базе данных", "error", err)
		os.Exit(1)
	}

	if err = urlRepo.Migrate(ctx, pool); err != nil {
		slog.Error("миграция базы данных", "error", err)
		os.Exit(1)
	}

	if err = run(cfg, log, pool); err != nil {
		slog.Error("ошибка запуска urlShortener сервиса", "error", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config, log *zap.Logger, pool *pgxpool.Pool) error {
	// Создать HTTP-обработчик через фабрику приложения.
	handler, err := app.NewHTTPHandler(cfg, log, pool)
	if err != nil {
		return fmt.Errorf("создание HTTP handler: %w", err)
	}

	httpServer := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	// Запускаем HTTP-сервер в отдельной горутине.
	go func() {
		log.Info("HTTP сервер запущен", zap.String("address", cfg.ServerAddr))

		if listenErr := httpServer.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			log.Error("ошибка HTTP сервера", zap.Error(listenErr))
		}
	}()

	// Ожидаем сигнал завершения.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("остановка urlShortener")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if shutdownErr := httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
		log.Error("ошибка остановки HTTP сервера", zap.Error(shutdownErr))
	}

	log.Info("urlShortener остановлен")

	return nil
}
