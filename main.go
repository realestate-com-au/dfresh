package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/docker/docker/pkg/term"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tagfush"

	app.Commands = []cli.Command{
		{
			Name:    "creds",
			Aliases: []string{"c"},
			Usage:   "show registry creds",
			Action: func(c *cli.Context) error {
				stdin, stdout, stderr := term.StdStreams()
				logrus.SetOutput(stderr)

				cli := command.NewDockerCli(stdin, stdout, stderr)
				opts := cliflags.NewClientOptions()
				cli.Initialize(opts)

				creds, error := cli.CredentialsStore("").Get("https://index.docker.io/v1/")
				if error != nil {
					panic(error)
				}
				fmt.Println("user:", creds.Username, "pass:", creds.Password, "auth:", creds.Auth)

				return nil
			},
		},
	}

	app.Run(os.Args)

}
