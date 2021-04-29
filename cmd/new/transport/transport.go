package transport

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "transport [options] <name>",
		Args:  cobra.MaximumNArgs(1),
		Short: "[experimental] Create a new ruby transport named <name> using given options",
		Long:  `[experimental] Create a new ruby transport named <name> using given options`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")

	return tmp
}
