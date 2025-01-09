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

	return nil
}
