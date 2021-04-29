package update

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	force     bool
	targetDir string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "update [flags]",
		Short: "Update a module that has been created by or converted for use by PDK",
		Long:  `Update a module that has been created by or converted for use by PDK`,
		RunE:  utils.ExecutePDKCommand,
	}

	tmp.Flags().BoolVar(&force, "force", false, "Update the module automatically, with no prompts")
	tmp.Flags().StringVar(&targetDir, "target-dir", "", "The target directory where you want PDK to write the packagee")

	return tmp
}
