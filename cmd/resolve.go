package cmd

import (
	"errors"
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
			ref, err := reference.ParseNormalizedNamed(args[0])
			if err != nil {
				return err
			}
			_, isDigested := ref.(reference.Digested)
			if isDigested {
				return errors.New("reference already has a digest")
			}
			canonicalRef, err := client.Resolve(ref)
			if err != nil {
				return err
			}
			fmt.Println(reference.FamiliarString(canonicalRef))
			return nil
		},
	}
}
