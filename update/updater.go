package update

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

type Updater struct {
	client       rego.Client
	reportWriter io.Writer
}

func NewUpdater(client rego.Client, reportWriter io.Writer) *Updater {
	return &Updater{client: client, reportWriter: reportWriter}
}

var refRegexp = regexp.MustCompile(reference.NameRegexp.String() + "(?::" + reference.TagRegexp.String() + ")?@" + reference.DigestRegexp.String() + "\\b")

func (u *Updater) UpdateRefsInFiles(paths []string) error {
	for _, path := range paths {
		err := u.UpdateRefsInFile(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Updater) UpdateRefsInFile(path string) error {
	buffer := new(bytes.Buffer)
	input, err := os.Open(path)
	if err != nil {
		return err
	}

	err = u.UpdateRefsInStream(path, input, buffer)
	if err != nil {
		return err
	}
	input.Close()

	return ioutil.WriteFile(path, buffer.Bytes(), 0666)
}

func (u *Updater) UpdateRefsInStream(streamName string, input io.Reader, output io.Writer) (err error) {
	scanner := bufio.NewScanner(input)
	line := 0
	for scanner.Scan() {
		line++
		context := streamName + ":" + strconv.Itoa(line)
		updated := refRegexp.ReplaceAllStringFunc(scanner.Text(), func(s string) string {
			return u.updateRef(s, context)
		})
		fmt.Fprintln(output, updated)
	}
	return scanner.Err()
}

func (u *Updater) updateRef(s string, context string) string {
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
	newRef, err := u.client.Resolve(ref)
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
	fmt.Fprintf(u.reportWriter, "%s: %s\n  was %s\n  now %s\n", context, nameAndTag, oldDigest, newDigest)
	return reference.FamiliarString(newRef)
}
