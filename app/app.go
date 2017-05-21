package app

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/docker/docker/pkg/term"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/registry"
)

var dockerCli *command.DockerCli

func Init() error {
	stdin, stdout, stderr := term.StdStreams()
	logrus.SetOutput(stderr)
	dockerCli := command.NewDockerCli(stdin, stdout, stderr)
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
	fmt.Println(ref)

	// Resolve the Repository name from fqn to RepositoryInfo
	repoInfo, err := registry.ParseRepositoryInfo(ref)
	if err != nil {
		return tags, err
	}
	fmt.Println(repoInfo)

	return []string{"hello"}, nil
}
