package cmd

import (
	"fmt"

	"github.com/mdub/dfresh/app"
	"github.com/spf13/cobra"
)

func newCredsCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "creds",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			server := "https://index.docker.io/v1/"
			if len(args) > 0 {
				server = args[0]
			}

			creds, err := app.GetAuthFor(server)
			if err != nil {
				return err
			}

			fmt.Println("user:", creds.Username, "pass:", creds.Password, "auth:", creds.Auth)

			return nil
		},
	}
}
