package main

import (
	"context"

	"git.front.kjuulh.io/kjuulh/orbis/internal/app"
	"git.front.kjuulh.io/kjuulh/orbis/internal/processes"
	"github.com/spf13/cobra"
)

func newRoot(app *app.App) *cobra.Command {
	logger := app.Logger()

	cmd := &cobra.Command{
		Use:   "orbis",
		Short: "Orbis is a data workflow scheduler for all your batch and real-time needs",

		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("starting orbis")

			ctx := context.Background()

			return processes.
				NewApp(logger).
				Add(app.Scheduler()).
				Add(app.Worker()).
				WithCtrlC().
				Execute(ctx)
		},
	}

	return cmd
}
