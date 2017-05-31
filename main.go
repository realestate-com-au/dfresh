package main

import (
	"os"

	"github.com/realestate-com-au/dfresh/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
