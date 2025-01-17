package app

import (
	"fmt"

	"git.front.kjuulh.io/kjuulh/orbis/internal/persistence"
	"git.front.kjuulh.io/kjuulh/orbis/internal/utilities"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Postgres = utilities.Singleton(func() (*pgxpool.Pool, error) {
	if err := persistence.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return persistence.NewConnection()
})
