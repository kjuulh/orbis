package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"git.front.kjuulh.io/kjuulh/orbis/internal/worker/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type workProcessor interface {
	ProcessNext(ctx context.Context, worker_id uuid.UUID) error
}

//go:generate sqlc generate

type Worker struct {
	workerID uuid.UUID

	db            *pgxpool.Pool
	workProcessor workProcessor
	logger        *slog.Logger

	capacity uint
}

func NewWorker(
	db *pgxpool.Pool,
	logger *slog.Logger,
	workProcessor workProcessor,
) *Worker {
	return &Worker{
		workerID:      uuid.New(),
		db:            db,
		workProcessor: workProcessor,
		logger:        logger,

		capacity: 50,
	}
}

func (w *Worker) Setup(ctx context.Context) error {
	repo := repositories.New(w.db)

	w.logger.InfoContext(ctx, "setting up worker", "worker_id", w.workerID)
	if err := repo.RegisterWorker(
		ctx,
		&repositories.RegisterWorkerParams{
			WorkerID: w.workerID,
			Capacity: int32(w.capacity),
		},
	); err != nil {
		return nil
	}

	return nil
}

type Workers struct {
	Instances []WorkerInstance
}

func (w *Workers) Capacity() uint {
	capacity := uint(0)

	for _, worker := range w.Instances {
		capacity += worker.Capacity
	}

	return capacity
}

type WorkerInstance struct {
	WorkerID uuid.UUID
	Capacity uint
}

func (w *Worker) GetWorkers(ctx context.Context) (*Workers, error) {
	repo := repositories.New(w.db)

	dbInstances, err := repo.GetWorkers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find workers: %w", err)
	}

	instances := make([]WorkerInstance, 0, len(dbInstances))
	for _, dbInstance := range dbInstances {
		instances = append(instances, WorkerInstance{
			WorkerID: dbInstance.WorkerID,
			Capacity: uint(dbInstance.Capacity),
		})
	}

	return &Workers{
		Instances: instances,
	}, nil
}

func (w *Worker) Start(ctx context.Context) error {
	heartBeatCtx, heartBeatCancel := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		errorCount := 0

		for {
			select {
			case <-heartBeatCtx.Done():
				return
			case <-ticker.C:
				if err := w.updateHeartBeat(heartBeatCtx); err != nil {
					if errorCount >= 5 {
						panic(fmt.Errorf("worker failed to register heartbeat for a long time, panicing..., err: %w", err))
					}
					errorCount += 1
				} else {
					errorCount = 0
				}
			}
		}
	}()

	defer func() {
		heartBeatCancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := w.processWorkQueue(ctx); err != nil {
				// FIXME: dead letter item, right now we just log and continue

				w.logger.WarnContext(ctx, "failed to handle work item", "error", err)
			}
		}
	}
}

func (w *Worker) updateHeartBeat(ctx context.Context) error {
	repo := repositories.New(w.db)

	w.logger.DebugContext(ctx, "updating heartbeat", "time", time.Now())
	return repo.UpdateWorkerHeartbeat(ctx, w.workerID)
}

func (w *Worker) processWorkQueue(ctx context.Context) error {
	return w.workProcessor.ProcessNext(ctx, w.workerID)
}
