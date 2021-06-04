package main

import (
	"github.com/puppetlabs/pdkgo/cmd/build"
	"github.com/puppetlabs/pdkgo/cmd/bundle"
	"github.com/puppetlabs/pdkgo/cmd/completion"
	"github.com/puppetlabs/pdkgo/cmd/console"
	"github.com/puppetlabs/pdkgo/cmd/convert"
	"github.com/puppetlabs/pdkgo/cmd/env"
	"github.com/puppetlabs/pdkgo/cmd/get"
	getConfig "github.com/puppetlabs/pdkgo/cmd/get/config"
	"github.com/puppetlabs/pdkgo/cmd/new"
	"github.com/puppetlabs/pdkgo/cmd/release"
	"github.com/puppetlabs/pdkgo/cmd/release/prep"
	"github.com/puppetlabs/pdkgo/cmd/release/publish"
	"github.com/puppetlabs/pdkgo/cmd/remove"
	removeConfig "github.com/puppetlabs/pdkgo/cmd/remove/config"
	"github.com/puppetlabs/pdkgo/cmd/root"
	"github.com/puppetlabs/pdkgo/cmd/set"
	setConfig "github.com/puppetlabs/pdkgo/cmd/set/config"
	"github.com/puppetlabs/pdkgo/cmd/test"
	"github.com/puppetlabs/pdkgo/cmd/test/unit"
	"github.com/puppetlabs/pdkgo/cmd/update"
	"github.com/puppetlabs/pdkgo/cmd/validate"
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

	testCmd := test.CreateCommand()
	testCmd.AddCommand(unit.CreateCommand())
	rootCmd.AddCommand(testCmd)

	rootCmd.AddCommand(build.CreateCommand())
	rootCmd.AddCommand(convert.CreateCommand())
	rootCmd.AddCommand(update.CreateCommand())

	releaseCmd := release.CreateCommand()
	releaseCmd.AddCommand(publish.CreateCommand())
	releaseCmd.AddCommand(prep.CreateCommand())
	rootCmd.AddCommand(releaseCmd)

	rootCmd.AddCommand(env.CreateCommand())
	rootCmd.AddCommand(validate.CreateCommand())

	setCmd := set.CreateCommand()
	setCmd.AddCommand(setConfig.CreateCommand())
	rootCmd.AddCommand(setCmd)

	getCmd := get.CreateCommand()
	getCmd.AddCommand(getConfig.CreateCommand())
	rootCmd.AddCommand(getCmd)

	removeCmd := remove.CreateCommand()
	removeCmd.AddCommand(removeConfig.CreateCommand())
	rootCmd.AddCommand(removeCmd)

	newCmd := new.CreateCommand()
	rootCmd.AddCommand(newCmd)

	rootCmd.AddCommand(bundle.CreateCommand())
	rootCmd.AddCommand(console.CreateCommand())

	cobra.OnInitialize(root.InitLogger, root.InitConfig)
	cobra.CheckErr(rootCmd.Execute())
}
