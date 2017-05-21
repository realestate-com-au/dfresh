package cmd

import (
	"fmt"

	"github.com/mdub/dfresh/app"

	"github.com/spf13/cobra"
)

func newTagsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tags IMAGE",
		Short: "Print all available tags for an image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tags, err := app.GetTags(args[0])
			if err != nil {
				return err
			}
			for _, tag := range tags {
				fmt.Println(tag)
			}
			return nil
		},
	}
}
