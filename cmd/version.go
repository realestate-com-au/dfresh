package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of dfresh",
		Long:  `All software has versions. This is dfresh's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("1.1.2")
		},
	}
}
