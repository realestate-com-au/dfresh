package cmd

import (
	"fmt"

	"github.com/mdub/dfresh/app"

	"github.com/spf13/cobra"
)

func newResolveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resolve IMAGE:TAG",
		Short: "Resolve an image reference",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			descriptor, err := app.Resolve(args[0])
			if err != nil {
				return err
			}
			fmt.Println(descriptor)
			return nil
		},
	}
}
