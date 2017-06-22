package cmd

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/realestate-com-au/dfresh/check"
	rego "github.com/realestate-com-au/dfresh/registry"
	"github.com/spf13/cobra"
)

func newCheckCmd(client rego.Client) *cobra.Command {

	var quiet bool

	command := &cobra.Command{
		Use:   "check [flags] [FILE...]",
		Short: "Check freshness of image references",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reportDestination := io.Writer(os.Stdout)
			if quiet {
				reportDestination = ioutil.Discard
			}
			checker := check.NewChecker(client, reportDestination)
			err := checker.CheckFiles(args, false)
			if err != nil {
				return err
			}
			if checker.UpdateCount() > 0 {
				os.Exit(1)
			}
			return nil
		},
	}
	command.Flags().BoolVarP(&quiet, "quiet", "q", false, "be silent")

	return command

}
