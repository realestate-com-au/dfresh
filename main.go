package main

import (
	"fmt"
	"os"

	"github.com/mdub/dfresh/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
