package app

import (
	"errors"

	"github.com/Sirupsen/logrus"
	dist "github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/docker/docker/distribution"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/docker/registry"
	"golang.org/x/net/context"
)

var dockerCli *command.DockerCli

// Initialise the app
func Init(debug bool) error {
	stdin, stdout, stderr := term.StdStreams()
	logrus.SetOutput(stderr)
	dockerCli = command.NewDockerCli(stdin, stdout, stderr)
	opts := cliflags.NewClientOptions()
	err := dockerCli.Initialize(opts)
	if err != nil {
		return err
	}
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled")
	}
	return nil
}

func newRepository(ctx context.Context, ref reference.Named) (dist.Repository, error) {

	repoInfo, err := registry.ParseRepositoryInfo(ref)
	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"ref":  ref,
		"repo": repoInfo,
	}).Debug("repository found")

	authConfig := command.ResolveAuthConfig(ctx, dockerCli, repoInfo.Index)
	registryService := registry.NewService(registry.ServiceOptions{V2Only: true})
	endpoints, err := registryService.LookupPullEndpoints(reference.Domain(repoInfo.Name))
	if err != nil {
		return nil, err
	}

	for _, endpoint := range endpoints {
		if endpoint.Version == registry.APIVersion1 {
			continue
		}
		repository, confirmedV2, err := distribution.NewV2Repository(ctx, repoInfo, endpoint, nil, &authConfig, "pull")
		if err == nil && confirmedV2 {
			return repository, nil
		}
	}

	return nil, errors.New("no V2 endpoint found")

}

// Get tags for a repository
func GetTags(s string) ([]string, error) {
	var tags []string
	ref, err := reference.ParseNormalizedNamed(s)
	if err != nil {
		return tags, err
	}
	if _, ok := ref.(reference.Tagged); ok {
		return tags, errors.New("reference already has a tag")
	}
	ctx := context.Background()
	repository, err := newRepository(ctx, ref)
	if err != nil {
		return tags, err
	}
	return repository.Tags(ctx).All(ctx)
}
