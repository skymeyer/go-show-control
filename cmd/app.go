package cmd

import (
	"github.com/spf13/cobra"

	// Load all drivers
	_ "go.skymyer.dev/show-control/dmx/driver/enttec"
	_ "go.skymyer.dev/show-control/dmx/driver/virtual"
)

func App() *cobra.Command {
	app := &cobra.Command{
		Use:   "showctl",
		Short: "Show Control",
	}

	app.AddCommand(
		NewShowCmd(),
		NewTestCmd(),
	)

	return app
}
