package cmd

import (
	"fmt"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/registry"
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

			// Resolve the Repository name from fqn to RepositoryInfo
			repoInfo, err := registry.ParseRepositoryInfo(ref)
			if err != nil {
				return err
			}
			fmt.Println(repoInfo)

			return nil
		},
	}
}
