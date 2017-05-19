package cmd

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/docker/docker/pkg/term"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stdin, stdout, stderr := term.StdStreams()
		logrus.SetOutput(stderr)

		cli := command.NewDockerCli(stdin, stdout, stderr)
		opts := cliflags.NewClientOptions()
		cli.Initialize(opts)

		creds, error := cli.CredentialsStore("").Get("https://index.docker.io/v1/")
		if error != nil {
			return error
		}
		fmt.Println("user:", creds.Username, "pass:", creds.Password, "auth:", creds.Auth)

		return nil
	},
}
