package unit

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	cleanFixtures          bool
	listUnitTestFiles      bool
	parallelUnitTests      bool
	puppetDevSourceVersion string
	puppetVersion          string
	unitTestsToRun         string
	verboseUnitTestOutput  bool
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "unit [flags]",
		Short: "Run unit tests",
		Long:  `Run unit tests`,
		RunE:  utils.ExecutePDKCommand,
	}

	tmp.Flags().BoolVarP(&cleanFixtures, "clean-fixtures", "c", false, "clean up downloaded fixtures after the test run")
	tmp.Flags().BoolVar(&listUnitTestFiles, "list", false, "list all available unit test files")
	tmp.Flags().BoolVar(&parallelUnitTests, "parallel", false, "run unit tests in parallel")

	tmp.Flags().StringVar(&puppetDevSourceVersion, "puppet-dev", "", "When specified, PDK will validate or test against the current Puppet source from github.com. To use this option, you must have network access to https://github.com.")
	tmp.Flags().StringVar(&puppetVersion, "puppet-version", "", "Puppet version to run tests or validations against")
	tmp.Flags().StringVar(&unitTestsToRun, "tests", "", "Specify a comma-separated list of unit test files to run. (default: )")

	tmp.Flags().BoolVar(&verboseUnitTestOutput, "verbose", false,
		"more verbose --list output. displays a list of examples in each unit test file")
	return tmp
}
