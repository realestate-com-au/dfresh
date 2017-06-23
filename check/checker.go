package check

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/reference"
	rego "github.com/realestate-com-au/dfresh/registry"
)

type Checker struct {
	client       rego.Client
	reportWriter io.Writer
	updateCount  int
}

func NewChecker(client rego.Client, reportWriter io.Writer) *Checker {
	return &Checker{client: client, reportWriter: reportWriter}
}

var refRegexp = regexp.MustCompile(reference.NameRegexp.String() + "(?::" + reference.TagRegexp.String() + ")?@" + reference.DigestRegexp.String() + "\\b")

func (c *Checker) CheckFiles(paths []string, saveUpdates bool) error {
	for _, path := range paths {
		err := c.CheckFile(path, saveUpdates)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Checker) CheckFile(path string, saveUpdates bool) error {
	buffer := new(bytes.Buffer)
	input, err := os.Open(path)
	if err != nil {
		return err
	}

	err = c.CheckStream(path, input, buffer)
	if err != nil {
		return err
	}
	input.Close()

	if saveUpdates {
		return ioutil.WriteFile(path, buffer.Bytes(), 0666)
	}
	return nil
}

func (c *Checker) UpdateCount() int {
	return c.updateCount
}

func (c *Checker) CheckStream(streamName string, input io.Reader, output io.Writer) (err error) {
	scanner := bufio.NewScanner(input)
	line := 0
	for scanner.Scan() {
		line++
		context := streamName + ":" + strconv.Itoa(line)
		updated := refRegexp.ReplaceAllStringFunc(scanner.Text(), func(s string) string {
			return c.updateRef(s, context)
		})
		fmt.Fprintln(output, updated)
	}
	return scanner.Err()
}

func (c *Checker) updateRef(s string, context string) string {
	parts := strings.Split(s, "@")
	nameAndTag := parts[0]
	oldDigest := parts[1]
	logrus.WithFields(logrus.Fields{
		"context":   context,
		"reference": nameAndTag,
	}).Debug("resolving")
	ref, err := reference.ParseNormalizedNamed(nameAndTag)
	if err != nil {
		panic(err)
	}
	newRef, err := c.client.Resolve(ref)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"reference": nameAndTag,
		}).Error("cannot resolve")
		return s
	}
	newDigest := newRef.Digest().String()
	if newDigest == oldDigest {
		return s
	}
	fmt.Fprintf(c.reportWriter, "%s: %s\n  old %s\n  new %s\n", context, nameAndTag, oldDigest, newDigest)
	c.updateCount++
	return reference.FamiliarString(newRef)
}
