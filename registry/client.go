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
	"github.com/docker/docker/api/types"
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

// This is mostly copied from the docker code base, and we deliberately don't
// ask the daemon for the default index server, since we don't want to rely on
// one running.
func (c *defaultClient) resolveAuthConfig(repoInfo *dreg.RepositoryInfo) types.AuthConfig {
	configKey := repoInfo.Index.Name
	if repoInfo.Index.Official {
		configKey = dreg.IndexServer
	}

	authConfig, _ := c.dockerCli.CredentialsStore(configKey).Get(configKey)
	return authConfig
}

func (c *defaultClient) newRepository(ref reference.Named) (dist.Repository, error) {

	repoInfo, err := dreg.ParseRepositoryInfo(ref)
	if err != nil {
		return nil, err
	}

	authConfig := c.resolveAuthConfig(repoInfo)

	registryService := dreg.NewService(dreg.ServiceOptions{V2Only: true})
	endpoints, err := registryService.LookupPullEndpoints(reference.Domain(repoInfo.Name))
	if err != nil {
		return nil, err
	}

	for _, endpoint := range endpoints {
		if endpoint.Version == dreg.APIVersion1 {
			continue
		}
		logrus.WithFields(logrus.Fields{
			"endpoint": endpoint.URL,
		}).Debug("contacting " + repoInfo.Index.Name)
		repository, confirmedV2, err := ddist.NewV2Repository(c.ctx, repoInfo, endpoint, nil, &authConfig, "pull")
		if err == nil && confirmedV2 {
			return repository, nil
		}
	}
	if err == nil {
		err = errors.New("cannot reach " + repoInfo.Index.Name)
	}

	return nil, err

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
