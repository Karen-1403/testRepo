package utils

import (
	"log/slog"
	"os"
)

// Logger is the global structured logger
var Logger *slog.Logger

// InitLogger initializes the global logger
func InitLogger() error {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	Logger = slog.New(handler)
	Logger.Info("logger initialized")

	return nil
}