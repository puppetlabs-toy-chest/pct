package build

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	force       bool
	targetDir   string
	noop        bool
	templateRef bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "build [flags]",
		Short: "Builds a package from the module that can be published to the Puppet Forge",
		Long:  `Builds a package from the module that can be published to the Puppet Forge`,
		RunE:  utils.ExecutePDKCommand,
	}

	tmp.Flags().BoolVar(&force, "force", false, "Skips the prompts and builds the module package")
	tmp.Flags().StringVar(&targetDir, "target-dir", "", "The target directory where you want PDK to write the package")
	tmp.Flags().BoolVar(&noop, "noop", false, "Do not update the module, just output what would be done")
	tmp.Flags().BoolVar(&templateRef, "templateRef", false, "Specifies the template git branch or tag to use when creating new modules or classes")

	return tmp
}
