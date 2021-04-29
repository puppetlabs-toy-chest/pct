package get

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "get [subcommand] [options]",
		Short: "Retrieve information about the PDK or current project",
		Long:  `Retrieve information about the PDK or current project`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")

	return tmp
}
