package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"go.skymyer.dev/show-control/show"
)

func NewShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Start Show Control",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			shutdownChan := make(chan os.Signal, 1)
			signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

			defer signal.Stop(shutdownChan)

			sc, err := show.NewFromConfig(
				"etc/setup1.yaml",
				show.WithFixtureLibrary("etc/fixtures.yaml"),
			)
			if err != nil {
				return err
			}
			sc.Run(ctx)

			<-shutdownChan
			return nil
		},
	}
}
