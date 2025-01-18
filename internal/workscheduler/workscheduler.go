package workscheduler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"git.front.kjuulh.io/kjuulh/orbis/internal/modelschedule"
	"git.front.kjuulh.io/kjuulh/orbis/internal/worker"
	"git.front.kjuulh.io/kjuulh/orbis/internal/workscheduler/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate sqlc generate

type WorkScheduler struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewWorkScheduler(
	db *pgxpool.Pool,
	logger *slog.Logger,
) *WorkScheduler {
	return &WorkScheduler{
		db:     db,
		logger: logger,
	}
}

type Worker struct {
	Instance          worker.WorkerInstance
	RemainingCapacity uint
}

type Workers struct {
	Workers []*Worker
}

func (w *Workers) IterateSlice(size uint) func(yield func([]Worker) bool) {
	return func(yield func([]Worker) bool) {
		if len(w.Workers) == 0 {
			return
		}

		workers := make([]Worker, 0)
		acc := uint(0)

		for {
			exit := true

			for _, worker := range w.Workers {
				if acc == size {
					if !yield(workers) {
						return
					}
					workers = make([]Worker, 0)
					acc = uint(0)

				}

				if worker.RemainingCapacity <= 0 {
					continue
				}

				worker.RemainingCapacity--
				workers = append(workers, *worker)
				acc++

				exit = false
			}

			if exit {
				if len(workers) > 0 {
					if !yield(workers) {
						return
					}
				}

				return
			}
		}
	}
}

func (w *WorkScheduler) GetWorkers(ctx context.Context, registeredWorkers *worker.Workers) (*Workers, error) {

	w.logger.DebugContext(ctx, "found workers", "workers", len(registeredWorkers.Instances))

	workers := make([]*Worker, 0, len(registeredWorkers.Instances))
	for _, registeredWorker := range registeredWorkers.Instances {
		remainingCapacity, err := w.GetWorker(ctx, &registeredWorker)
		if err != nil {
			return nil, fmt.Errorf("failed to find capacity for worker: %w", err)
		}

		if remainingCapacity == 0 {
			w.logger.DebugContext(ctx, "skipping worker as no remaining capacity")
			continue
		}

		workers = append(workers, &Worker{
			Instance:          registeredWorker,
			RemainingCapacity: remainingCapacity,
		})
	}

	return &Workers{Workers: workers}, nil
}

func (w *WorkScheduler) GetWorker(
	ctx context.Context,
	worker *worker.WorkerInstance,
) (uint, error) {
	repo := repositories.New(w.db)

	current_size, err := repo.GetCurrentQueueSize(ctx, worker.WorkerID)
	if err != nil {
		return 0, fmt.Errorf("failed to get current queue size: %s: %w", worker.WorkerID, err)
	}

	if int64(worker.Capacity)-current_size <= 0 {
		return 0, nil
	}

	return worker.Capacity - uint(current_size), nil
}

func (w *WorkScheduler) InsertModelRun(
	ctx context.Context,
	worker Worker,
	modelRun *modelschedule.ModelRunSchedule,
) error {
	repo := repositories.New(w.db)

	return repo.InsertQueueItem(ctx, &repositories.InsertQueueItemParams{
		ScheduleID: uuid.New(),
		WorkerID:   worker.Instance.WorkerID,
		StartRun: pgtype.Timestamptz{
			Time:  modelRun.Start,
			Valid: true,
		},
		EndRun: pgtype.Timestamptz{
			Time:  modelRun.End,
			Valid: true,
		},
	})
}

func (w *WorkScheduler) GetNext(ctx context.Context, workerID uuid.UUID) (*uuid.UUID, error) {
	repo := repositories.New(w.db)

	schedule, err := repo.GetNext(ctx, workerID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get next worker item: %w", err)
		}

		return nil, nil
	}

	return &schedule.ScheduleID, nil
}

func (w *WorkScheduler) StartProcessing(ctx context.Context, scheduleID uuid.UUID) error {
	repo := repositories.New(w.db)

	return repo.StartProcessing(ctx, scheduleID)
}

func (w *WorkScheduler) Archive(ctx context.Context, scheduleID uuid.UUID) error {
	repo := repositories.New(w.db)

	return repo.Archive(ctx, scheduleID)
}
