package utils

import (
	"log/slog"
	"os"
)

var globalLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func NewLogger(env string) *slog.Logger {
	level := slog.LevelInfo
	if env == "local" {
		level = slog.LevelDebug
	}

	handlerOptions := &slog.HandlerOptions{Level: level}

	if env == "local" {
		return slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
}

func SetLogger(logger *slog.Logger) {
	if logger != nil {
		globalLogger = logger
	}
}

func Logger() *slog.Logger {
	return globalLogger
}
