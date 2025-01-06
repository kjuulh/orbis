package scheduler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
)

type Scheduler struct {
	logger *slog.Logger
	db     *pgx.Conn
}

func NewScheduler(logger *slog.Logger, db *pgx.Conn) *Scheduler {
	return &Scheduler{
		logger: logger,
		db:     db,
	}
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

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("gracefully shutting down elected scheduler")
			return nil
		default:
			if err := s.process(ctx); err != nil {
				return fmt.Errorf("scheduler failed: %w", err)
			}
		}
	}
}

func (s *Scheduler) acquireLeader(ctx context.Context) (bool, error) {
	for {
		select {
		case <-ctx.Done():
			return false, nil

		default:
			var acquiredLock bool
			if err := s.db.QueryRow(ctx, "SELECT pg_try_advisory_lock(1234)").Scan(&acquiredLock); err != nil {
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
	s.logger.Debug("scheduler processing items")

	// FIXME: simulate work
	time.Sleep(time.Second * 2)

	return nil
}
