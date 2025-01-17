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

//go:generate sqlc generate

type Worker struct {
	workerID uuid.UUID

	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewWorker(
	db *pgxpool.Pool,
	logger *slog.Logger,
) *Worker {
	return &Worker{
		workerID: uuid.New(),
		db:       db,
		logger:   logger,
	}
}

func (w *Worker) Setup(ctx context.Context) error {
	repo := repositories.New(w.db)

	if err := repo.RegisterWorker(ctx, w.workerID); err != nil {
		return nil
	}

	return nil
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
		if err := w.processWorkQueue(ctx); err != nil {
			// FIXME: dead letter item, right now we just log and continue

			w.logger.WarnContext(ctx, "failed to handle work item", "error", err)
		}
	}
}

func (w *Worker) updateHeartBeat(ctx context.Context) error {
	repo := repositories.New(w.db)

	w.logger.DebugContext(ctx, "updating heartbeat", "time", time.Now())
	return repo.UpdateWorkerHeartbeat(ctx, w.workerID)
}

func (w *Worker) processWorkQueue(_ context.Context) error {
	time.Sleep(time.Second)

	return nil
}
