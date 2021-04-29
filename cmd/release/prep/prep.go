package prep

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug             bool
	force             bool
	file              string
	forgeToken        string
	forgeUploadUrl    string
	skipBuild         bool
	skipChangelog     bool
	skipDependecy     bool
	skipDocumentation bool
	skipPublish       bool
	skipValidation    bool
	versionBump       string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "prep [flags]",
		Short: "(Experimental) Performs all the pre-release checks to ensure module is ready to be released",
		Long:  `(Experimental) Performs all the pre-release checks to ensure module is ready to be released`,
		RunE:  utils.ExecutePDKCommand,
	}

	tmp.Flags().BoolVar(&force, "force", false, "Release the module automatically, with no prompts")
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	tmp.Flags().StringVar(&forgeToken, "forge-token", "", "Set Forge API token")
	tmp.Flags().StringVar(&forgeUploadUrl, "forge-upload-url", "", "Set forge upload url path. (default: https://forgeapi.puppetlabs.com/v3/releases)")

	tmp.Flags().StringVar(&file, "file", "", "Path to the built module to push to the Forge. This option can only be used when --skip-build is also used. Defaults to pkg/<module version>.tar.gz")

	tmp.Flags().BoolVar(&skipBuild, "skip-build", false, "Skips module build")
	tmp.Flags().BoolVar(&skipChangelog, "skip-changelog", false, "Skips the automatic changelog generation")

	tmp.Flags().BoolVar(&skipDependecy, "skip-dependency", false, "Skips the module dependency check")
	tmp.Flags().BoolVar(&skipDocumentation, "skip-documentation", false, "Skips the documentation update")
	tmp.Flags().BoolVar(&skipPublish, "skip-publish", false, "Skips publishing the module to the forge")

	tmp.Flags().BoolVar(&skipValidation, "skip-validation", false, "Skips the module validation check")

	tmp.Flags().StringVar(&versionBump, "version", "", "Update the module to the specified version prior to release. When not specified, the new version will be computed from the Changelog where possible")
	return tmp
}
