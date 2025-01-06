package main

import "github.com/spf13/cobra"

func newRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orbis",
		Short: "Orbis is a data workflow scheduler for all your batch and real-time needs",
	}

	return cmd
}
