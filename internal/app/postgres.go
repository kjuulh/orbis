package app

import (
	"fmt"

	"git.front.kjuulh.io/kjuulh/orbis/internal/persistence"
	"git.front.kjuulh.io/kjuulh/orbis/internal/utilities"
	"github.com/jackc/pgx/v5"
)

var Postgres = utilities.Singleton(func() (*pgx.Conn, error) {
	if err := persistence.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return persistence.NewConnection()
})
