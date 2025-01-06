package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"git.front.kjuulh.io/kjuulh/orbis/internal/utilities"
	"github.com/jackc/pgx/v5"
)

var Postgres = utilities.Singleton(func() (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, err := pgx.Connect(ctx, os.Getenv("ORBIS_POSTGRES_DB"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to orbis postgres database: %w", err)
	}

	return conn, nil
})
