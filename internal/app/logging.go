package app

import (
	"log/slog"
	"os"
	"strings"

	"gitlab.com/greyxor/slogor"
)

func setupLogging() *slog.Logger {
	logLevelRaw := os.Getenv("ORBIS_LOG_LEVEL")
	var logLevel slog.Leveler
	switch strings.ToLower(logLevelRaw) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	return slog.New(slogor.NewHandler(os.Stderr, slogor.SetLevel(logLevel)))
}
