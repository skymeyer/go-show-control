package cmd

import (
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/common"
	"go.skymyer.dev/show-control/config"

	// Load dmx and io drivers
	_ "go.skymyer.dev/show-control/dmx/driver/artnet"
	_ "go.skymyer.dev/show-control/dmx/driver/enttec"
	_ "go.skymyer.dev/show-control/dmx/driver/virtual"
	_ "go.skymyer.dev/show-control/io/console"
	_ "go.skymyer.dev/show-control/io/novation"
	_ "go.skymyer.dev/show-control/io/webhook"
)

var cfg *config.App

func App() *cobra.Command {
	var (
		debug   bool
		stats   time.Duration
		cfgFile = "etc/app.yaml"
	)
	app := &cobra.Command{
		Use:   "showctl",
		Short: "Show Control",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// Load app configuration
			cfg = &config.App{}
			if err := common.LoadFromFile(cfgFile, cfg); err != nil {
				return err
			}

			// Logger setup
			var (
				zcfg zap.Config
				err  error
			)

			if debug {
				zcfg = zap.NewDevelopmentConfig()
			} else {
				zcfg = zap.NewProductionConfig()
			}

			// Use log file
			if cfg.LogFile != "" {
				zcfg.OutputPaths = []string{cfg.LogFile}
			}
			logger.Default, err = zcfg.Build()
			if err != nil {
				return err
			}

			// Enable periodic global system stats
			if stats > 0 {
				withSystemStats(stats)
			}

			logger.Default.Info("runtime init", zap.String("log", cfg.LogFile), zap.Bool("debug", debug))
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return logger.Default.Sync()
		},
	}

	app.AddCommand(
		NewRunCmd(),
	)

	app.PersistentFlags().BoolVar(&debug, "debug", debug, "Debug logging")
	app.PersistentFlags().DurationVar(&stats, "stats", stats, "Dump debug system stats window")
	app.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "Configuration file")

	return app
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
