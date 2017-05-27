package registry

import (
	"errors"
	"os"

	"github.com/Sirupsen/logrus"
	dist "github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	clicmd "github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	ddist "github.com/docker/docker/distribution"
	dreg "github.com/docker/docker/registry"
	"golang.org/x/net/context"
)

type Client interface {
	Init(debug bool) error
	GetTags(ref reference.Named) ([]string, error)
	Resolve(ref reference.Named) (reference.Canonical, error)
}

type defaultClient struct {
	ctx       context.Context
	dockerCli *clicmd.DockerCli
}

func NewClient() Client {
	return &defaultClient{ctx: context.Background()}
}

// Initialise the app
func (c *defaultClient) Init(debug bool) error {
	c.dockerCli = clicmd.NewDockerCli(os.Stdin, os.Stdout, os.Stderr)
	opts := cliflags.NewClientOptions()
	opts.Common.Debug = debug
	return c.dockerCli.Initialize(opts)
}

func (c *defaultClient) newRepository(ref reference.Named) (dist.Repository, error) {

	repoInfo, err := dreg.ParseRepositoryInfo(ref)
	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"ref":  ref,
		"repo": repoInfo,
	}).Debug("repository found")

	authConfig := clicmd.ResolveAuthConfig(c.ctx, c.dockerCli, repoInfo.Index)
	registryService := dreg.NewService(dreg.ServiceOptions{V2Only: true})
	endpoints, err := registryService.LookupPullEndpoints(reference.Domain(repoInfo.Name))
	if err != nil {
		return nil, err
	}

	for _, endpoint := range endpoints {
		if endpoint.Version == dreg.APIVersion1 {
			continue
		}
		repository, confirmedV2, err := ddist.NewV2Repository(c.ctx, repoInfo, endpoint, nil, &authConfig, "pull")
		if err == nil && confirmedV2 {
			return repository, nil
		}
	}

	return nil, errors.New("no V2 endpoint found")

}

func (c *defaultClient) GetTags(ref reference.Named) ([]string, error) {
	var tags []string
	repository, err := c.newRepository(ref)
	if err != nil {
		return tags, err
	}
	return repository.Tags(c.ctx).All(c.ctx)
}

func (c *defaultClient) Resolve(ref reference.Named) (reference.Canonical, error) {
	canonicalRef, isCanonical := ref.(reference.Canonical)
	if isCanonical {
		return canonicalRef, nil
	}
	tag := "latest"
	taggedRef, isTagged := ref.(reference.Tagged)
	if isTagged {
		tag = taggedRef.Tag()
	}
	repository, err := c.newRepository(ref)
	if err != nil {
		return nil, err
	}
	descriptor, err := repository.Tags(c.ctx).Get(c.ctx, tag)
	if err != nil {
		return nil, err
	}
	return reference.WithDigest(ref, descriptor.Digest)
}
