package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaximVebDevJs/go-u-shr/internal/config"
	"github.com/MaximVebDevJs/go-u-shr/internal/logger"
	"github.com/MaximVebDevJs/go-u-shr/pkg/app"
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

	if err = run(cfg, log); err != nil {
		slog.Error("ошибка запуска urlShortener сервиса", "error", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config, log *zap.Logger) error {
	// Создать HTTP-обработчик через фабрику приложения.
	handler := app.NewHTTPHandler(cfg.BaseURL, log)

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
