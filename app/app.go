package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/docker/docker/pkg/term"
)

var dockerCli *command.DockerCli

func Init() error {
	stdin, stdout, stderr := term.StdStreams()
	logrus.SetOutput(stderr)
	cli := command.NewDockerCli(stdin, stdout, stderr)
	opts := cliflags.NewClientOptions()
	err := cli.Initialize(opts)
	if err != nil {
		return err
	}
	dockerCli = cli
	return nil
}

func GetAuthFor(server string) (types.AuthConfig, error) {
	return dockerCli.CredentialsStore("").Get(server)
}
