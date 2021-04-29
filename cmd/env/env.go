package env

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug         bool
	force         bool
	puppetDev     bool
	peVersion     string
	puppetVersion string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "env [flags]",
		Short: "Aids in setting a CLI context for a specified version of Puppet by outputting export commands for necessary environment variables",
		Long:  `Aids in setting a CLI context for a specified version of Puppet by outputting export commands for necessary environment variables`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&force, "force", false, "Release the module automatically, with no prompts")
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	tmp.Flags().StringVar(&peVersion, "pe-version", "", "Puppet Enterprise version to run tests or validations against")
	tmp.Flags().StringVar(&puppetVersion, "puppet-version", "", "Puppet Enterprise version to run tests or validations against")
	tmp.Flags().BoolVar(&puppetDev, "puppet-dev", false, "When specified, PDK will validate or test against the current Puppet source from github.com. To use this option, you must have network access to https://github.com")

	return tmp
}
