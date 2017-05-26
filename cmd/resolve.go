package cmd

import (
	"fmt"

	rego "github.com/mdub/dfresh/registry"

	"github.com/spf13/cobra"
)

func newResolveCmd(client rego.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve IMAGE:TAG",
		Short: "Resolve an image reference",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			digest, err := client.GetDigest(args[0])
			if err != nil {
				return err
			}
			fmt.Println(digest)
			return nil
		},
	}
}
