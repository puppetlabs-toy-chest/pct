package main

import (
	"github.com/puppetlabs/pdkgo/cmd/pct/completion"
	"github.com/puppetlabs/pdkgo/cmd/pmr/get"
	"github.com/puppetlabs/pdkgo/cmd/pmr/root"
	"github.com/puppetlabs/pdkgo/cmd/pmr/set"
	appver "github.com/puppetlabs/pdkgo/cmd/pmr/version"

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

	rootCmd.AddCommand(set.CreateSetCommand())
	rootCmd.AddCommand(get.CreateGetCommand())

	rootCmd.AddCommand(completion.CreateCompletionCommand())

	// afero setup
	// fs := afero.NewOsFs()
	// afs := afero.Afero{Fs: fs}
	// iofs := afero.IOFS{Fs: fs}

	cobra.OnInitialize(root.InitLogger, root.InitConfig)
	cobra.CheckErr(rootCmd.Execute())
}
