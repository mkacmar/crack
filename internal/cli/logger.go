package cli

import (
	"io"
	"log/slog"
)

const logLevelDisabled = slog.Level(99)

func setupLogger(level string, output io.Writer) *slog.Logger {
	var logLevel slog.Level

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
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(output, opts)
	return slog.New(handler)
}
