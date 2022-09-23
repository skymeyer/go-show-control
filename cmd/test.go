package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"go.skymyer.dev/show-control/config"
	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/dmx/driver"
)

func NewTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Test stuff",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			shutdownChan := make(chan os.Signal, 1)
			signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

			defer signal.Stop(shutdownChan)

			data, err := ioutil.ReadFile("etc/setup1.yaml")
			if err != nil {
				return err
			}

			var cfg config.Setup
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return err
			}

			var f = dmx.NewDMX512Frame()

			dp, err := driver.New(
				cfg.Devices["dmxpro1"].Driver,
				cfg.Devices["dmxpro1"].Device,
			)
			if err != nil {
				return err
			}
			dp.SetUniverse(0, f)
			defer dp.Close()

			dp.Run(ctx)

			go func() {
				ticker := time.NewTicker(300 * time.Millisecond)
				defer ticker.Stop()

				f.SetSlot(22, 0xFF)
				for {
					select {
					case <-ctx.Done():
						fmt.Printf("\nstopped slot writer\n")
						return
					case <-ticker.C:
						f.SetSlot(23, 0xFF)
						f.SetSlot(24, 0x00)
						f.SetSlot(25, 0x00)
						time.Sleep(100 * time.Millisecond)

						f.SetSlot(23, 0x00)
						f.SetSlot(24, 0xFF)
						f.SetSlot(25, 0x00)
						time.Sleep(100 * time.Millisecond)

						f.SetSlot(23, 0x00)
						f.SetSlot(24, 0x00)
						f.SetSlot(25, 0xFF)
						time.Sleep(100 * time.Millisecond)
					}
				}
			}()

			fmt.Println("running ...")
			<-shutdownChan
			fmt.Printf("\nbye\n")
			return nil
		},
	}
}
