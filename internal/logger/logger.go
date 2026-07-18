package logger

import (
	"go.uber.org/zap"
)

func New(level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	cfg.Level = lvl

	return cfg.Build()
}

// Nop заглушка для тестов
func Nop() *zap.Logger {
	return zap.NewNop()
}
