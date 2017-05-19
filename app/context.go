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

var DefaultContext = &Context{}

func (context *Context) Init() error {
	stdin, stdout, stderr := term.StdStreams()
	logrus.SetOutput(stderr)
	cli := command.NewDockerCli(stdin, stdout, stderr)
	opts := cliflags.NewClientOptions()
	err := cli.Initialize(opts)
	if err != nil {
		return err
	}
	context.dockerCli = cli
	return nil
}

func (context *Context) GetAuthFor(server string) (types.AuthConfig, error) {
	return context.dockerCli.CredentialsStore("").Get(server)
}
