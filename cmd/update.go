package cmd

import (
	"io"
	"io/ioutil"
	"os"

	rego "github.com/realestate-com-au/dfresh/registry"
	"github.com/realestate-com-au/dfresh/update"
	"github.com/spf13/cobra"
)

func newUpdateCmd(client rego.Client) *cobra.Command {

	var quiet bool

	command := &cobra.Command{
		Use:   "update [flags] [FILE...]",
		Short: "Update image references",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			reportDestination := io.Writer(os.Stdout)
			if quiet || len(args) == 0 {
				reportDestination = ioutil.Discard
			}
			updater := update.NewUpdater(client, reportDestination)
			if len(args) == 0 {
				return updater.UpdateRefsInStream("-", os.Stdin, os.Stdout)
			}
			return updater.UpdateRefsInFiles(args)
		},
	}
	command.Flags().BoolVarP(&quiet, "quiet", "q", false, "be silent")

	return command

}
