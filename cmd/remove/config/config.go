package config

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug bool
	force bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "config [name] [value]",
		Short: "Remove or delete the configuration for <name>",
		Long:  `Remove or delete the configuration for <name>`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	tmp.Flags().BoolVar(&force, "force", false, "Force multi-value configuration settings to be removed instead of emptied")

	return tmp
}
