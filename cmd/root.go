package cmd

import (
	"github.com/mdub/dfresh/app"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {

	var debug bool

	root := &cobra.Command{
		Use: "dfresh",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return app.Init(debug)
		},
	}

	root.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debugging")

	root.AddCommand(newTagsCmd())

	return root

}
