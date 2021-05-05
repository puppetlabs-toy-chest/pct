package class

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "class [options] <class_name>",
		Args:  cobra.MaximumNArgs(1),
		Short: "Create a new class named <class_name> using given options",
		Long:  `Create a new class named <class_name> using given options`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")

	return tmp
}
