package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newTagsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tags IMAGE",
		Short: "Print all available tags for an image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			name := args[0]
			fmt.Println(name)

			return nil
		},
	}
}
