package cmd

import (
	"fmt"

	"github.com/docker/distribution/reference"
	rego "github.com/mdub/dfresh/registry"
	"github.com/spf13/cobra"
)

func newResolveCmd(client rego.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve IMAGE:TAG",
		Short: "Resolve an image reference",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			canonicalRef, err := client.Resolve(args[0])
			if err != nil {
				return err
			}
			fmt.Println(reference.FamiliarString(canonicalRef))
			return nil
		},
	}
}
