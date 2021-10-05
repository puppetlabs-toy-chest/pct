package version

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/puppetlabs/pdkgo/pkg/telemetry"
	"github.com/spf13/cobra"
)

func CreateVersionCommand(version, buildDate string, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			_, span := telemetry.NewSpan(cmd.Context(), "version")
			defer telemetry.EndSpan(span)
			telemetry.AddStringSpanAttribute(span, "name", "version")

			fmt.Fprint(os.Stdout, Format(version, buildDate, commit))
		},
	}

	return cmd
}

func Format(version, buildDate string, commit string) string {
	version = strings.TrimSpace(strings.TrimPrefix(version, "v"))

	var dateStr string
	if buildDate != "" {
		t, _ := time.Parse(time.RFC3339, buildDate)
		dateStr = t.Format("2006/01/02")
	}

	if commit != "" && len(commit) > 7 {
		length := len(commit) - 7
		commit = strings.TrimSpace(commit[:len(commit)-length])
	}

	return fmt.Sprintf("pct %s %s %s\n\n%s",
		version, commit, dateStr, changelogURL(version))
}

func changelogURL(version string) string {
	path := "https://github.com/puppetlabs/pdkgo"
	r := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[\w.]+)?$`)
	if !r.MatchString(version) {
		return fmt.Sprintf("%s/releases/latest", path)
	}

	url := fmt.Sprintf("%s/releases/tag/%s", path, strings.TrimPrefix(version, "v"))
	return url
}
