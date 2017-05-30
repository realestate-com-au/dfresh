package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/reference"
	rego "github.com/mdub/dfresh/registry"
	"github.com/spf13/cobra"
)

var refRegexp = regexp.MustCompile(reference.NameRegexp.String() + "(?::" + reference.TagRegexp.String() + ")?@" + reference.DigestRegexp.String() + "\\b")

func newUpdateCmd(client rego.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "update [flags] [FILE...]",
		Short: "Update image references",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return updateRefsInStream(client, os.Stdin, os.Stdout)
			} else {
				return updateRefsInFiles(client, args)
			}
		},
	}
}

func updateRefsInFiles(client rego.Client, paths []string) error {
	for _, path := range paths {
		err := updateRefsInFile(client, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateRefsInFile(client rego.Client, path string) error {
	buffer := new(bytes.Buffer)
	input, err := os.Open(path)
	if err != nil {
		return err
	}

	err = updateRefsInStream(client, input, buffer)
	if err != nil {
		return err
	}
	input.Close()

	return ioutil.WriteFile(path, buffer.Bytes(), 0666)
}

func updateRefsInStream(client rego.Client, input io.Reader, output io.Writer) (err error) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		updated := refRegexp.ReplaceAllStringFunc(scanner.Text(), func(s string) string {
			return getUpdatedRef(client, s)
		})
		fmt.Fprintln(output, updated)
	}
	return scanner.Err()
}

func getUpdatedRef(client rego.Client, s string) string {
	parts := strings.Split(s, "@")
	nameAndTag := parts[0]
	oldDigest := parts[1]
	ref, err := reference.ParseNormalizedNamed(nameAndTag)
	if err != nil {
		panic(err)
	}
	newRef, err := client.Resolve(ref)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"reference": nameAndTag,
		}).Error("cannot resolve")
		return s
	}
	newDigest := newRef.Digest().String()
	if newDigest != oldDigest {
		logrus.WithFields(logrus.Fields{
			"was": oldDigest,
			"now": newDigest,
		}).Debug("updated " + nameAndTag)
	}
	return reference.FamiliarString(newRef)
}
