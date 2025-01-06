package app

import (
	"log/slog"
)

type App struct {
	logger *slog.Logger
}

func NewApp() *App {
	return &App{
		logger: setupLogging(),
	}
}

func (a *App) Logger() *slog.Logger {
	return a.logger
}
