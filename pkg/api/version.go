package api

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func GetVersionString() string {
	return format(version, date, commit)
}

func format(version, buildDate string, commit string) string {
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
