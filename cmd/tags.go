package cmd

import (
	"fmt"

	rego "github.com/mdub/dfresh/registry"
	"github.com/spf13/cobra"
)

func newTagsCmd(client rego.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "tags IMAGE",
		Short: "Print all available tags for an image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tags, err := client.GetTags(args[0])
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
