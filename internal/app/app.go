package app

import (
	"log/slog"

	"git.front.kjuulh.io/kjuulh/orbis/internal/executor"
	"git.front.kjuulh.io/kjuulh/orbis/internal/modelschedule"
	"git.front.kjuulh.io/kjuulh/orbis/internal/scheduler"
	"git.front.kjuulh.io/kjuulh/orbis/internal/worker"
	"git.front.kjuulh.io/kjuulh/orbis/internal/workprocessor"
	"git.front.kjuulh.io/kjuulh/orbis/internal/workscheduler"
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
	return scheduler.NewScheduler(a.logger.With("component", "scheduler"), Postgres(), a.Executor(), a.Worker())
}

func (a *App) Executor() *executor.Executor {
	return executor.NewExecutor(
		a.logger.With("component", "executor"),
		ModelRegistry(),
		a.ModelSchedule(),
		a.Worker(),
		a.WorkScheduler(),
	)
}

func (a *App) Worker() *worker.Worker {
	return worker.NewWorker(Postgres(), a.logger, a.WorkProcessor())
}

func (a *App) WorkScheduler() *workscheduler.WorkScheduler {
	return workscheduler.NewWorkScheduler(Postgres(), a.logger)
}

func (a *App) WorkProcessor() *workprocessor.WorkProcessor {
	return workprocessor.NewWorkProcessor(a.WorkScheduler(), a.logger)
}

func (a *App) ModelSchedule() *modelschedule.ModelSchedule {
	return modelschedule.NewModelSchedule(a.logger, Postgres())
}
