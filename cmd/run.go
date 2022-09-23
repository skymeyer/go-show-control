package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"go.skymyer.dev/show-control/app"
	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/common"
)

func NewRunCmd() *cobra.Command {
	var (
		debug bool
		reset bool
		stats time.Duration
	)
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run show control",
		RunE: func(cmd *cobra.Command, args []string) error {

			if debug {
				logger.Default, _ = zap.NewDevelopment()
				if stats > 0 {
					withSystemStats(stats)
				}
			}
			defer logger.Default.Sync()

			controller, err := app.New()
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

	cmd.Flags().BoolVar(&debug, "debug", debug, "Debug logging")
	cmd.Flags().BoolVar(&reset, "reset", reset, "Reset database")
	cmd.Flags().DurationVar(&stats, "stats", stats, "Dump debug system stats window")

	return cmd
}

func withSystemStats(each time.Duration) {
	ticker := time.NewTicker(each)
	go func() {
		defer ticker.Stop()
		for t := range ticker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			logger.Default.Debug("system stats",
				zap.Time("ticker", t),
				zap.Uint64("alloc", m.Alloc),
				zap.Uint64("alloc_total", m.TotalAlloc),
				zap.Uint64("sys", m.Sys),
				zap.Uint64("gc", uint64(m.NumGC)),
				zap.Uint64("gc_forced", uint64(m.NumForcedGC)),
			)
		}
	}()
}
