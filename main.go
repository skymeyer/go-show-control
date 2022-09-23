package main

import (
	"os"

	"go.skymyer.dev/show-control/cmd"
)

func main() {
	if err := cmd.App().Execute(); err != nil {
		os.Exit(1)
	}
}
