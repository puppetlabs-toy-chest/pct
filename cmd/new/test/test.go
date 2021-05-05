package test

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug         bool
	unit          bool
	puppetDev     bool
	peVersion     string
	puppetVersion string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "test [options] <name>",
		Args:  cobra.MaximumNArgs(1),
		Short: "Create a new test for the object named <name>",
		Long:  `Create a new test for the object named <name>`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	tmp.Flags().BoolVar(&unit, "unit", false, "Create a new unit test")

	tmp.Flags().StringVar(&peVersion, "pe-version", "", "Puppet Enterprise version to run tests or validations against")
	tmp.Flags().StringVar(&puppetVersion, "puppet-version", "", "Puppet Enterprise version to run tests or validations against")
	tmp.Flags().BoolVar(&puppetDev, "puppet-dev", false, "When specified, PDK will validate or test against the current Puppet source from github.com. To use this option, you must have network access to https://github.com")

	return tmp
}
