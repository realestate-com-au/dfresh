package cmd

import (
	"fmt"

	"github.com/mdub/tagfush/app"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:  "tagfush",
	Args: cobra.MaximumNArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return app.DefaultContext.Init()
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		server := "https://index.docker.io/v1/"
		if len(args) > 0 {
			server = args[0]
		}

		creds, err := app.DefaultContext.GetAuthFor(server)
		if err != nil {
			return err
		}

		fmt.Println("user:", creds.Username, "pass:", creds.Password, "auth:", creds.Auth)

		return nil
	},
}
