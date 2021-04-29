package set

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "set [subcommand] [options]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Set or update information about the PDK or current project",
		Long:  `Set or update information about the PDK or current project`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")

	return tmp
}
