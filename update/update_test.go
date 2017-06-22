package update

import (
	"bytes"
	"github.com/docker/distribution/reference"
	"github.com/opencontainers/go-digest"
	"strings"
	"testing"
)

type stubClient struct {
	digest digest.Digest
}

func (c *stubClient) Init(debug bool) error {
	return nil
}

func (c *stubClient) GetTags(ref reference.Named) ([]string, error) {
	return []string{}, nil
}

func (c *stubClient) Resolve(name reference.Named) (reference.Canonical, error) {
	return reference.WithDigest(name, c.digest)
}

func TestUpdateRefsInStream(t *testing.T) {
	client := &stubClient{digest: "sha256:08868d719684cf9cafacbaa1786ad01111332b4c1e65abd67833db603d8dab7f"}
	inputReader := strings.NewReader("ruby:2.3@sha256:a5ebd3bc0bf3881258975f8afa1c6d24429dfd4d7dd53a299559a3e927b77fd7")
	outputWriter := new(bytes.Buffer)

	err := UpdateRefsInStream(client, inputReader, outputWriter)

	if err != nil {
		t.Error("Did not expect error, ", err)
	}

	expectedOutput := "ruby:2.3@sha256:08868d719684cf9cafacbaa1786ad01111332b4c1e65abd67833db603d8dab7f\n"
	output := outputWriter.String()
	if output != expectedOutput {
		t.Error(
			"expected", expectedOutput,
			"got", output,
		)
	}
}
