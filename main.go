package main

import (
	"github.com/puppetlabs/pdkgo/cmd/completion"
	"github.com/puppetlabs/pdkgo/cmd/root"
	"github.com/puppetlabs/pdkgo/cmd/test"
	"github.com/puppetlabs/pdkgo/cmd/test/unit"
	appver "github.com/puppetlabs/pdkgo/cmd/version"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var rootCmd = root.CreateRootCommand()

	var verCmd = appver.CreateVersionCommand(version, date, commit)
	v := appver.Format(version, date, commit)
	rootCmd.Version = v
	rootCmd.SetVersionTemplate(v)
	rootCmd.AddCommand(verCmd)

	rootCmd.AddCommand(completion.CreateCompletionCommand())

	testCmd := test.CreateTestCommand()
	testCmd.AddCommand(unit.CreateTestUnitCommand())
	rootCmd.AddCommand(testCmd)

	cobra.OnInitialize(root.InitConfig)
	cobra.CheckErr(rootCmd.Execute())
}
