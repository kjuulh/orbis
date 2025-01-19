package deadletter

import (
	"context"
	"fmt"
	"log/slog"

	"git.front.kjuulh.io/kjuulh/orbis/internal/deadletter/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeadLetter struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewDeadLetter(db *pgxpool.Pool, logger *slog.Logger) *DeadLetter {
	return &DeadLetter{
		db:     db,
		logger: logger,
	}
}

func (d *DeadLetter) InsertDeadLetter(ctx context.Context, schedule uuid.UUID) error {
	repo := repositories.New(d.db)

	d.logger.WarnContext(ctx, "deadlettering schedule", "schedule", schedule)
	if err := repo.InsertDeadLetter(ctx, schedule); err != nil {
		return fmt.Errorf("failed to insert item into dead letter: %w", err)
	}

	return nil

}
