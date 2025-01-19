package workprocessor

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"git.front.kjuulh.io/kjuulh/orbis/internal/deadletter"
	"git.front.kjuulh.io/kjuulh/orbis/internal/workscheduler"
	"github.com/google/uuid"
)

type WorkProcessor struct {
	workscheduler *workscheduler.WorkScheduler
	logger        *slog.Logger
	deadletter    *deadletter.DeadLetter
}

func NewWorkProcessor(
	workscheduler *workscheduler.WorkScheduler,
	logger *slog.Logger,
	deadletter *deadletter.DeadLetter,
) *WorkProcessor {
	return &WorkProcessor{
		workscheduler: workscheduler,
		logger:        logger,
		deadletter:    deadletter,
	}
}

func (w *WorkProcessor) ProcessNext(ctx context.Context, workerID uuid.UUID) error {
	schedule, err := w.workscheduler.GetNext(ctx, workerID)
	if err != nil {
		return fmt.Errorf("failed to get next work item: %w", err)
	}

	if schedule == nil {
		// TODO: defer somewhere else
		time.Sleep(time.Second)
		return nil
	}

	w.logger.DebugContext(ctx, "handling item", "schedule", schedule)

	if err := w.workscheduler.StartProcessing(ctx, *schedule); err != nil {
		return fmt.Errorf("failed to start processing items: %w", err)
	}

	time.Sleep(time.Millisecond * 10)
	// Process or achive

	if err := w.workscheduler.Archive(ctx, *schedule); err != nil {
		return fmt.Errorf("failed to archive item: %w", err)
	}

	if rand.Float64() > 0.9 {
		if err := w.deadletter.InsertDeadLetter(ctx, *schedule); err != nil {
			return fmt.Errorf("failed to deadletter message: %w", err)
		}
	}

	return nil
}
