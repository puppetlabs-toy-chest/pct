package main

import (
	"github.com/puppetlabs/pdkgo/pkg/cmd/completion"
	"github.com/puppetlabs/pdkgo/pkg/cmd/root"
	appver "github.com/puppetlabs/pdkgo/pkg/cmd/version"
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

	cobra.OnInitialize(root.InitConfig)
	cobra.CheckErr(rootCmd.Execute())
}
