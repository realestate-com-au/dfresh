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
					newRef, err := client.Resolve(parts[0])
					if err != nil {
						logrus.Error(fmt.Sprintf("while resolving %s: %s", parts[0], err))
						return s
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
