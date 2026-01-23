package cli

import (
	"log/slog"
	"os"
)

const logLevelDisabled = slog.Level(99)

func setupLogger(level string) (*slog.Logger, bool) {
	var logLevel slog.Level
	valid := true

	switch level {
	case "none":
		logLevel = logLevelDisabled
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelError
		valid = false
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(os.Stderr, opts)
	return slog.New(handler), valid
}
