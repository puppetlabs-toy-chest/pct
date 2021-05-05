package convert

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	force           bool
	addTests        bool
	defaultTemplate bool
	fullInterview   bool
	noop            bool
	skipInterview   bool
	templateRef     bool
	templateUrl     bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "convert [flags]",
		Short: "Convert an existing module to be compatible with the PDK",
		Long:  `Convert an existing module to be compatible with the PDK`,
		RunE:  utils.ExecutePDKCommand,
	}

	tmp.Flags().BoolVar(&force, "force", false, "Convert the module automatically, with no prompts")
	tmp.Flags().BoolVar(&addTests, "addTests", false, "Add any missing tests while converting the module")
	tmp.Flags().BoolVar(&defaultTemplate, "defaultTemplate", false, "Convert the module to use the default PDK template")
	tmp.Flags().BoolVar(&fullInterview, "fullInterview", false, "When specified, interactive querying of metadata will include all optional questions")
	tmp.Flags().BoolVar(&noop, "noop", false, "Do not convert the module, just output what would be done")
	tmp.Flags().BoolVar(&skipInterview, "skipInterview", false, "When specified, skips interactive querying of metadata")
	tmp.Flags().BoolVar(&templateRef, "templateRef", false, "Specifies the template git branch or tag to use when creating new modules or classes")
	tmp.Flags().BoolVar(&templateUrl, "templateUrl", false, "Specifies the URL to the template to use when creating new modules or classes. (default: pdk-default#2.1.0)")

	return tmp
}
