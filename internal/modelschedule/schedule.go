package modelschedule

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"git.front.kjuulh.io/kjuulh/orbis/internal/modelregistry"
	"git.front.kjuulh.io/kjuulh/orbis/internal/modelschedule/repositories"
	"github.com/adhocore/gronx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate sqlc generate

type ModelRunSchedule struct {
	Model *modelregistry.Model
	Start time.Time
	End   time.Time
}

type ModelSchedule struct {
	logger *slog.Logger

	db *pgxpool.Pool
}

func NewModelSchedule(logger *slog.Logger, db *pgxpool.Pool) *ModelSchedule {
	return &ModelSchedule{
		logger: logger,

		db: db,
	}
}

func (m *ModelSchedule) GetNext(
	ctx context.Context,
	model modelregistry.Model,
	start time.Time,
	end time.Time,
	amount uint,
) (models []ModelRunSchedule, lastExecuted *time.Time, err error) {
	repo := repositories.New(m.db)

	var startRun time.Time
	lastRun, err := repo.GetLast(ctx, model.Name)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, fmt.Errorf("failed to get last run for mode: %s: %w", model.Name, err)
		}

		startRun = start
	} else {
		startRun = lastRun.Time
	}

	times := make([]ModelRunSchedule, 0, amount)
	for {
		next, err := gronx.NextTickAfter(model.Schedule, startRun, false)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to find next model schedule: %w", err)
		}

		if next.Equal(time.Time{}) {
			break
		}

		if next.After(end) {
			break
		}

		times = append(times, ModelRunSchedule{
			Model: &model,
			Start: startRun,
			End:   next,
		})
		startRun = next

		if len(times) >= int(amount) {
			break
		}
	}

	if len(times) == 0 {
		return nil, nil, nil
	}

	return times, &startRun, nil
}

func (m *ModelSchedule) UpdateModelRun(ctx context.Context, model modelregistry.Model, lastRun *time.Time) error {
	repo := repositories.New(m.db)

	return repo.UpsertModel(ctx, &repositories.UpsertModelParams{
		ModelName: model.Name,
		LastRun: pgtype.Timestamptz{
			Time:  *lastRun,
			Valid: true,
		},
	})
}
