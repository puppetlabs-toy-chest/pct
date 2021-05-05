package config

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "config [name]",
		Short: "Retrieve the configuration for <name>. If not specified, retrieve all configuration settings",
		Long:  `Retrieve the configuration for <name>. If not specified, retrieve all configuration settings`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")

	return tmp
}
