package main

import (
	dreg "github.com/docker/docker/registry"
	"github.com/realestate-com-au/dfresh/cmd"
	"os"
)

func main() {
	dreg.CertsDir = ""
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(10)
	}
}
