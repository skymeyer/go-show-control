package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"go.skymyer.dev/show-control/app"
	"go.skymyer.dev/show-control/common"
)

func NewRunCmd() *cobra.Command {
	var (
		reset bool
	)
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run show control",
		RunE: func(cmd *cobra.Command, args []string) error {

			controller, err := app.New(cfg)
			if err != nil {
				return err
			}

			// FIXME: proper controller config
			if reset {
				os.Remove(filepath.Join(common.MustUserConfigDir(app.Name), app.Database))
			}

			return controller.Run()
		},
	}

	cmd.Flags().BoolVar(&reset, "reset", reset, "Reset database")

	return cmd
}
