package executor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"git.front.kjuulh.io/kjuulh/orbis/internal/modelregistry"
	"git.front.kjuulh.io/kjuulh/orbis/internal/modelschedule"
	"git.front.kjuulh.io/kjuulh/orbis/internal/worker"
	"git.front.kjuulh.io/kjuulh/orbis/internal/workscheduler"
)

type Executor struct {
	logger *slog.Logger

	modelRegistry   *modelregistry.ModelRegistry
	modelSchedule   *modelschedule.ModelSchedule
	worker          *worker.Worker
	workerscheduler *workscheduler.WorkScheduler
}

func NewExecutor(
	logger *slog.Logger,
	modelRegistry *modelregistry.ModelRegistry,
	modelSchedule *modelschedule.ModelSchedule,
	worker *worker.Worker,
	workerscheduler *workscheduler.WorkScheduler,
) *Executor {
	return &Executor{
		logger: logger,

		modelRegistry:   modelRegistry,
		modelSchedule:   modelSchedule,
		worker:          worker,
		workerscheduler: workerscheduler,
	}
}

func (e *Executor) DispatchEvents(ctx context.Context) error {

	e.logger.InfoContext(ctx, "dispatching events")

	start := time.Now().Add(-time.Second * 30)
	end := time.Now()

	models, err := e.modelRegistry.GetModels()
	if err != nil {
		return fmt.Errorf("failed to get models from registry: %w", err)
	}

	registeredWorkers, err := e.worker.GetWorkers(ctx)
	if err != nil {
		return fmt.Errorf("failed to find workers: %w", err)
	}

	e.logger.InfoContext(ctx, "moving unattended events")
	if err := e.workerscheduler.GetUnattended(ctx, registeredWorkers); err != nil {
		return fmt.Errorf("failed to move unattended events: %w", err)
	}

	workers, err := e.workerscheduler.GetWorkers(ctx, registeredWorkers)
	if err != nil {
		return fmt.Errorf("failed to find workers: %w", err)
	}

	for workers := range workers.IterateSlice(2000) {
		for _, model := range models {
			modelRuns, lastRun, err := e.modelSchedule.GetNext(ctx, model, start, end, uint(len(workers)))
			if err != nil {
				return err
			}

			for i, modelRun := range modelRuns {
				worker := workers[i]
				e.logger.DebugContext(ctx, "dispatching model run", "modelRun", modelRun.Model.Name, "start", modelRun.Start, "end", modelRun.End)

				if err := e.workerscheduler.InsertModelRun(ctx, worker, &modelRun); err != nil {
					return fmt.Errorf("failed to register model run: %w", err)
				}
			}

			if err := e.modelSchedule.UpdateModelRun(ctx, model, lastRun); err != nil {
				return fmt.Errorf("failed to update checkpoint for model: %w", err)
			}
		}
	}

	return nil
}
