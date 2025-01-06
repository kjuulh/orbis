package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Scheduler struct {
	logger *slog.Logger
}

func NewScheduler(logger *slog.Logger) *Scheduler {
	return &Scheduler{
		logger: logger,
	}
}

func (s *Scheduler) Execute(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("gracefully shutting down scheduler")
			return nil
		default:
			if err := s.process(ctx); err != nil {
				return fmt.Errorf("scheduler failed: %w", err)
			}
		}
	}
}

func (s *Scheduler) process(ctx context.Context) error {
	s.logger.Debug("scheduler processing items")

	// FIXME: simulate work
	time.Sleep(time.Second * 5)

	return nil
}
