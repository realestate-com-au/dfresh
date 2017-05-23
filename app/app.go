package app

import (
	"github.com/Sirupsen/logrus"
	dist "github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/docker/docker/distribution"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/docker/registry"
	"golang.org/x/net/context"
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

	registryService := registry.NewService(registry.ServiceOptions{V2Only: true})
	_, _, err = registryService.Auth(ctx, &authConfig, "dfresh")
	if err != nil {
		return tags, err
	}

	// get endpoints
	endpoints, err := registryService.LookupPullEndpoints(reference.Domain(repoInfo.Name))
	if err != nil {
		return tags, err
	}

	// retrieve repository
	var (
		confirmedV2 bool
		repository  dist.Repository
		lastError   error
	)

	for _, endpoint := range endpoints {
		if endpoint.Version == registry.APIVersion1 {
			continue
		}

		repository, confirmedV2, lastError = distribution.NewV2Repository(ctx, repoInfo, endpoint, nil, authConfig, "pull")
		if lastError == nil && confirmedV2 {
			break
		}
	}
	if lastError != nil && confirmedV2 {
		break
	}
	if lastError != nil {
		return tags, lastError
	}

	return repository.Tags(ctx).All(ctx)
}
