package app

import (
	"log/slog"
	"os"
)

func setGlobalLogger(level string) {
	opt := &slog.HandlerOptions{
		Level: parseLogLevel(level),
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opt)))
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "info":
		return slog.LevelInfo
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	default:
		return slog.LevelDebug
	}
}
