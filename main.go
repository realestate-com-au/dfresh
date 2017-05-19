package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/cli/command"
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
				fmt.Println(cli)
				return nil
			},
		},
	}

	app.Run(os.Args)

}
