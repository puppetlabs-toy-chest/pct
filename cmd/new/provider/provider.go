package provider

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "provider [options] <name>",
		Args:  cobra.MaximumNArgs(1),
		Short: "Create a new custom provider named <name> using given options",
		Long:  `Create a new custom provider named <name> using given options`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")

	return tmp
}
