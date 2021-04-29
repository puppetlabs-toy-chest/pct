package validate

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug         bool
	autoCorrect   bool
	list          bool
	parallel      bool
	puppetDev     bool
	peVersion     string
	puppetVersion string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "validate [validators] [options] [targets]",
		Short: "Run static analysis tests",
		Long: `Run metadata, YAML, Puppet, Ruby, or Tasks validation.

    [validators] is an optional comma-separated list of validators to use. If
    not specified, all validators are used. Note that when using PowerShell,
    the list of validators must be enclosed in single quotes.

    [targets] is an optional space-separated list of files or directories to
    be validated. If not specified, validators are run against all applicable
    files in the module`,
		RunE: utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	tmp.Flags().BoolVarP(&autoCorrect, "auto-correct", "a", false, "Automatically correct problems where possible")
	tmp.Flags().BoolVar(&list, "list", false, "List all available validators")
	tmp.Flags().BoolVar(&parallel, "parallel", false, "Run validations in parallel")

	tmp.Flags().StringVar(&peVersion, "pe-version", "", "Puppet Enterprise version to run tests or validations against")
	tmp.Flags().StringVar(&puppetVersion, "puppet-version", "", "Puppet Enterprise version to run tests or validations against")
	tmp.Flags().BoolVar(&puppetDev, "puppet-dev", false, "When specified, PDK will validate or test against the current Puppet source from github.com. To use this option, you must have network access to https://github.com")

	return tmp
}
