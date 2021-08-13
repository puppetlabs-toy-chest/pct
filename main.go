package main

import (
	"github.com/puppetlabs/pdkgo/cmd/build"
	"github.com/puppetlabs/pdkgo/cmd/completion"
	"github.com/puppetlabs/pdkgo/cmd/new"
	"github.com/puppetlabs/pdkgo/cmd/root"
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

	rootCmd.AddCommand(build.CreateCommand())
	rootCmd.AddCommand(new.CreateCommand())

	cobra.OnInitialize(root.InitLogger, root.InitConfig)
	cobra.CheckErr(rootCmd.Execute())
}
