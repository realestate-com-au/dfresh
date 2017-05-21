package cmd

import (
	"fmt"

	"github.com/docker/distribution/reference"
	"github.com/spf13/cobra"
)

func newTagsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tags IMAGE",
		Short: "Print all available tags for an image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			ref, err := reference.ParseNormalizedNamed(args[0])
			if err != nil {
				return err
			}

			fmt.Println(ref)

			return nil
		},
	}
}
