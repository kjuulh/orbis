package main

import (
	"fmt"

	"git.front.kjuulh.io/kjuulh/orbis/internal/app"
	"github.com/spf13/cobra"
)

func newRoot(app *app.App) *cobra.Command {
	logger := app.Logger()

	cmd := &cobra.Command{
		Use:   "orbis",
		Short: "Orbis is a data workflow scheduler for all your batch and real-time needs",

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger.Info("starting orbis")

			if err := app.Scheduler().Execute(ctx); err != nil {
				return fmt.Errorf("scheduler failed with error: %w", err)
			}

			return nil
		},
	}

	return cmd
}
