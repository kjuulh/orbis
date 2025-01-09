package app

import (
	"log/slog"

	"git.front.kjuulh.io/kjuulh/orbis/internal/executor"
	"git.front.kjuulh.io/kjuulh/orbis/internal/scheduler"
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

func (a *App) Scheduler() *scheduler.Scheduler {
	return scheduler.NewScheduler(a.logger.With("component", "scheduler"), Postgres(), a.Executor())
}

func (a *App) Executor() *executor.Executor {
	return executor.NewExecutor(a.logger.With("component", "executor"))
}
