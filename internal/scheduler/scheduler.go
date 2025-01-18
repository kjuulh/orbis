package scheduler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"git.front.kjuulh.io/kjuulh/orbis/internal/executor"
	"git.front.kjuulh.io/kjuulh/orbis/internal/worker"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Scheduler struct {
	logger   *slog.Logger
	db       *pgxpool.Pool
	executor *executor.Executor
	worker   *worker.Worker
}

func NewScheduler(logger *slog.Logger, db *pgxpool.Pool, executor *executor.Executor, worker *worker.Worker) *Scheduler {
	return &Scheduler{
		logger:   logger,
		db:       db,
		executor: executor,
		worker:   worker,
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	if err := s.Execute(ctx); err != nil {
		return fmt.Errorf("execution of scheduler failed: %w", err)
	}

	return nil
}

func (s *Scheduler) Execute(ctx context.Context) error {
	acquiredLeader, err := s.acquireLeader(ctx)
	if err != nil {
		return err
	}

	if !acquiredLeader {
		s.logger.Info("gracefully shutting down non-elected scheduler")
		return nil
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("gracefully shutting down elected scheduler")
			return nil
		case <-ticker.C:
			if err := s.process(ctx); err != nil {
				return fmt.Errorf("scheduler failed: %w", err)
			}
		}
	}
}

func (s *Scheduler) acquireLeader(ctx context.Context) (bool, error) {
	db, err := s.db.Acquire(ctx)

	for {
		select {
		case <-ctx.Done():
			return false, nil
		default:
			if err != nil {
				return false, fmt.Errorf("failed to acquire db connection: %w", err)
			}

			var acquiredLock bool
			if err := db.QueryRow(ctx, "SELECT pg_try_advisory_lock(1234)").Scan(&acquiredLock); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return false, nil
				}
			}

			if !acquiredLock {
				wait := time.Second * time.Duration(rand.Float32()*9+1)

				s.logger.Debug("failed to acquire lock, parking non-elected scheduler", "wait_seconds", wait)
				time.Sleep(wait)
				continue
			}

			s.logger.Info("acquired lock, electing application to leader")
			return true, nil

		}
	}
}

func (s *Scheduler) process(ctx context.Context) error {
	if err := s.worker.Prune(ctx); err != nil {
		return fmt.Errorf("failed to prune error: %w", err)
	}

	if err := s.executor.DispatchEvents(ctx); err != nil {
		return fmt.Errorf("failed to dispatch events: %w", err)
	}

	return nil
}
