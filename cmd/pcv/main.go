package main

import (
	"github.com/puppetlabs/pdkgo/cmd/pcv/build"
	"github.com/puppetlabs/pdkgo/cmd/pcv/root"
	"github.com/puppetlabs/pdkgo/cmd/pcv/run"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = root.CreateRootCommand()

	rootCmd.AddCommand(build.CreateBuildCommand())
	rootCmd.AddCommand(run.CreateRunCommand())

	cobra.CheckErr(rootCmd.Execute())
}
