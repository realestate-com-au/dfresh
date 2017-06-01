package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	rego "github.com/realestate-com-au/dfresh/registry"
	"github.com/spf13/cobra"
)

func initLogging(debug bool) {
	logrus.SetOutput(os.Stderr)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func NewRootCmd() *cobra.Command {

	var debug bool
	client := rego.NewClient()

	root := &cobra.Command{
		Use: "dfresh",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			initLogging(debug)
			return client.Init(debug)
		},
		SilenceUsage: true,
	}

	root.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debugging")

	root.AddCommand(newTagsCmd(client))
	root.AddCommand(newResolveCmd(client))
	root.AddCommand(newUpdateCmd(client))
	root.AddCommand(newVersionCmd())

	return root

}
