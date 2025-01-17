package persistence

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection() (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, err := pgxpool.New(ctx, os.Getenv("ORBIS_POSTGRES_DB"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to orbis postgres database: %w", err)
	}

	return conn, nil
}
