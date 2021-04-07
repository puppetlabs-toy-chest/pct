package version

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func CreateVersionCommand(version, buildDate string, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(os.Stdout, Format(version, buildDate, commit))
		},
	}

	return cmd
}

func Format(version, buildDate string, commit string) string {
	version = strings.TrimPrefix(version, "v")

	var dateStr string
	if buildDate != "" {
		dateStr = fmt.Sprintf(" (%s)", buildDate)
	}

	return fmt.Sprintf("pdk %s%s\npdk-ruby 2.2.0\n\n%s", version, dateStr, changelogURL(version))
}

func changelogURL(version string) string {
	path := "https://github.com/puppetlabs/pdkgo"
	r := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[\w.]+)?$`)
	if !r.MatchString(version) {
		return fmt.Sprintf("%s/releases/latest", path)
	}

	url := fmt.Sprintf("%s/releases/tag/v%s", path, strings.TrimPrefix(version, "v"))
	return url
}
