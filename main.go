package main

import (
	"os"

	"github.com/mdub/dfresh/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
