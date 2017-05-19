package cmd

import (
	"fmt"

	"github.com/mdub/tagfush/app"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:  "tagfush",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		context, err := app.NewContext()
		if err != nil {
			return err
		}

		server := "https://index.docker.io/v1/"
		if len(args) > 0 {
			server = args[0]
		}

		creds, err := context.GetAuthFor(server)
		fmt.Println("user:", creds.Username, "pass:", creds.Password, "auth:", creds.Auth)

		return nil
	},
}
