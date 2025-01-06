package app

import (
	"log/slog"
	"os"

	"gitlab.com/greyxor/slogor"
)

func setupLogging() *slog.Logger {
	return slog.New(slogor.NewHandler(os.Stderr))
}
