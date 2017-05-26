package registry

import (
	"errors"
	"os"

	"github.com/Sirupsen/logrus"
	dist "github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/docker/docker/distribution"
	"github.com/docker/docker/registry"
	"github.com/opencontainers/go-digest"
	"golang.org/x/net/context"
)

type Client interface {
	Init(debug bool) error
	GetTags(s string) ([]string, error)
	GetDigest(s string) (digest.Digest, error)
}

type defaultClient struct {
	ctx       context.Context
	dockerCli *command.DockerCli
}

func NewClient() Client {
	return &defaultClient{ctx: context.Background()}
}

// Initialise the app
func (c *defaultClient) Init(debug bool) error {
	c.dockerCli = command.NewDockerCli(os.Stdin, os.Stdout, os.Stderr)
	opts := cliflags.NewClientOptions()
	opts.Common.Debug = debug
	return c.dockerCli.Initialize(opts)
}

func (c *defaultClient) newRepository(ref reference.Named) (dist.Repository, error) {

	repoInfo, err := registry.ParseRepositoryInfo(ref)
	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"ref":  ref,
		"repo": repoInfo,
	}).Debug("repository found")

	authConfig := command.ResolveAuthConfig(c.ctx, c.dockerCli, repoInfo.Index)
	registryService := registry.NewService(registry.ServiceOptions{V2Only: true})
	endpoints, err := registryService.LookupPullEndpoints(reference.Domain(repoInfo.Name))
	if err != nil {
		return nil, err
	}

	for _, endpoint := range endpoints {
		if endpoint.Version == registry.APIVersion1 {
			continue
		}
		repository, confirmedV2, err := distribution.NewV2Repository(c.ctx, repoInfo, endpoint, nil, &authConfig, "pull")
		if err == nil && confirmedV2 {
			return repository, nil
		}
	}

	return nil, errors.New("no V2 endpoint found")

}

func (c *defaultClient) GetTags(s string) ([]string, error) {
	var tags []string
	ref, err := reference.ParseNormalizedNamed(s)
	if err != nil {
		return tags, err
	}
	if _, ok := ref.(reference.Tagged); ok {
		return tags, errors.New("reference already has a tag")
	}
	ctx := c.ctx
	repository, err := c.newRepository(ref)
	if err != nil {
		return tags, err
	}
	return repository.Tags(ctx).All(ctx)
}

func (c *defaultClient) GetDigest(s string) (digest.Digest, error) {
	var nullDigest digest.Digest
	ref, err := reference.ParseNormalizedNamed(s)
	if err != nil {
		return nullDigest, err
	}
	digestedRef, ok := ref.(reference.Digested)
	if ok {
		return digestedRef.Digest(), nil
	}
	tag := "latest"
	taggedRef, ok := ref.(reference.Tagged)
	if ok {
		tag = taggedRef.Tag()
	}
	repository, err := c.newRepository(ref)
	if err != nil {
		return nullDigest, err
	}
	ctx := c.ctx
	descriptor, err := repository.Tags(ctx).Get(ctx, tag)
	if err != nil {
		return nullDigest, err
	}
	return descriptor.Digest, nil
}
