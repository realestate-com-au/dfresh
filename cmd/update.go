package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/reference"
	rego "github.com/mdub/dfresh/registry"
	"github.com/spf13/cobra"
)

var refRegexp = regexp.MustCompile(reference.NameRegexp.String() + "(?::" + reference.TagRegexp.String() + ")?@" + reference.DigestRegexp.String())

func newUpdateCmd(client rego.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update image references",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				var updateRef = func(s string) string {
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
						}).Info("updated " + nameAndTag)
					}
					return reference.FamiliarString(newRef)
				}
				updated := refRegexp.ReplaceAllStringFunc(scanner.Text(), updateRef)
				fmt.Println(updated)
			}
			return scanner.Err()
		},
	}
}
