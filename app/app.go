package app

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/docker/docker/pkg/term"
	"golang.org/x/net/context"

	"github.com/docker/docker/registry"
)

var dockerCli *command.DockerCli

func Init() error {
	stdin, stdout, stderr := term.StdStreams()
	logrus.SetOutput(stderr)
	dockerCli = command.NewDockerCli(stdin, stdout, stderr)
	opts := cliflags.NewClientOptions()
	err := dockerCli.Initialize(opts)
	if err != nil {
		return err
	}
	return nil
}

func GetAuthFor(server string) (types.AuthConfig, error) {
	return dockerCli.CredentialsStore("").Get(server)
}

func GetTags(s string) ([]string, error) {
	var tags []string
	ref, err := reference.ParseNormalizedNamed(s)
	if err != nil {
		return tags, err
	}
	repoInfo, err := registry.ParseRepositoryInfo(ref)
	if err != nil {
		return tags, err
	}

	ctx := context.Background()

	authConfig := command.ResolveAuthConfig(ctx, dockerCli, repoInfo.Index)

	service := registry.NewService(registry.ServiceOptions{V2Only: true})
	_, _, err = service.Auth(ctx, &authConfig, "dfresh")
	if err != nil {
		return tags, err
	}

	return []string{"hello"}, nil
}
