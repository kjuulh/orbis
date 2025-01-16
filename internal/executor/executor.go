package executor

import (
	"context"
	"log/slog"
)

type Executor struct {
	logger *slog.Logger
}

func NewExecutor(logger *slog.Logger) *Executor {
	return &Executor{
		logger: logger,
	}
}

func (e *Executor) DispatchEvents(ctx context.Context) error {
	e.logger.InfoContext(ctx, "dispatching events")

	// TODO: Process updates to models
	// TODO: Insert new cron for runtime
	// TODO: Calculate time since last run
	// TODO: Send events for workers to pick up

	return nil
}
