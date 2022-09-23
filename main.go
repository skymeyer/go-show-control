package main

//go:generate go-bindata -prefix db/migrations -o db/migrations.go -pkg db -nomemcopy db/migrations/...

import (
	"os"

	"go.skymyer.dev/show-control/cmd"
)

func main() {
	if err := cmd.App().Execute(); err != nil {
		os.Exit(1)
	}
}
