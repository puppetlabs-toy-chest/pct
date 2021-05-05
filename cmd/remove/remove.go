package remove

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "remove [command]",
		Short: "Remove or delete information about the PDK or current project",
		Long:  `Remove or delete information about the PDK or current project`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")

	return tmp
}
