package cmd

import (
	"fmt"
	"os"

	"github.com/puppetlabs/pdkgo/pkg/api"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:    "version",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprint(os.Stdout, api.GetVersionString())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
