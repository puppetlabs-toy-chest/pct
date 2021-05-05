package definedtype

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "defined_type [options] <name>",
		Args:  cobra.MaximumNArgs(1),
		Short: "Create a new defined type named <name> using given options",
		Long:  `Create a new defined type named <name> using given options`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")

	return tmp
}
