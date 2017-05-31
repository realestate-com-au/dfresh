package cmd

import (
	"errors"
	"fmt"

	"github.com/docker/distribution/reference"
	rego "github.com/realestate-com-au/dfresh/registry"
	"github.com/spf13/cobra"
)

func newTagsCmd(client rego.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "tags IMAGE",
		Short: "Print all available tags for an image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := reference.ParseNormalizedNamed(args[0])
			if err != nil {
				return err
			}
			if _, ok := ref.(reference.Tagged); ok {
				return errors.New("reference already has a tag")
			}
			tags, err := client.GetTags(ref)
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
