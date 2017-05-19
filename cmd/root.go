package cmd

import (
	"github.com/mdub/dfresh/app"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {

	root := &cobra.Command{
		Use: "dfresh",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return app.DefaultContext.Init()
		},
	}
	root.AddCommand(newCredsCmd())
	return root

}
