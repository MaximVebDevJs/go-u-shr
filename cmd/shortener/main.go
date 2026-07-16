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
	"github.com/MaximVebDevJs/go-u-shr/pkg/app"
)

const (
	shutdownTimeout = 10 * time.Second
	readTimeout     = 15 * time.Second
	writeTimeout    = 15 * time.Second
	idleTimeout     = 60 * time.Second
)

func main() {
	config.ParseFlags()

	if err := run(); err != nil {
		slog.Error("ошибка запуска urlShortener сервиса", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Создать HTTP-обработчик через фабрику приложения.
	handler := app.NewHTTPHandler()

	httpServer := &http.Server{
		Addr:         config.FlagRunAddr,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	// Запускаем HTTP-сервер в отдельной горутине.
	go func() {
		slog.Info("HTTP сервер запущен", "address", config.FlagRunAddr)

		if listenErr := httpServer.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			slog.Error("ошибка HTTP сервера", "error", listenErr)
		}
	}()

	// Ожидаем сигнал завершения.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("остановка urlShortener")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if shutdownErr := httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
		slog.Error("ошибка остановки HTTP сервера", "error", shutdownErr)
	}

	slog.Info("urlShortener остановлен")

	return nil
}
