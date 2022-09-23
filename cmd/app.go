package cmd

import (
	"github.com/spf13/cobra"

	// Load dmx and midi device support
	_ "go.skymyer.dev/show-control/dmx/driver/enttec"
	_ "go.skymyer.dev/show-control/dmx/driver/virtual"
	_ "go.skymyer.dev/show-control/io/novation"
)

func App() *cobra.Command {
	app := &cobra.Command{
		Use:   "showctl",
		Short: "Show Control",
	}

	app.AddCommand(
		NewRunCmd(),
	)

	return app
}
