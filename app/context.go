package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/docker/docker/pkg/term"
)

type Context struct {
	dockerCli *command.DockerCli
}

func NewContext() (*Context, error) {
	stdin, stdout, stderr := term.StdStreams()
	logrus.SetOutput(stderr)

	cli := command.NewDockerCli(stdin, stdout, stderr)
	opts := cliflags.NewClientOptions()
	err := cli.Initialize(opts)
	if err != nil {
		return nil, err
	}

	return &Context{dockerCli: cli}, nil
}

func (context *Context) GetAuthFor(server string) (types.AuthConfig, error) {
	return context.dockerCli.CredentialsStore("").Get(server)
}
