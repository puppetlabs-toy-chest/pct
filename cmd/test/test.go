package test

import (
	"github.com/spf13/cobra"
)

var (
	debug  bool
	format string
)

func CreateTestCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "test [flags]",
		Short: "Run tests",
		Long:  `Run tests`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	tmp.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug output")
	tmp.PersistentFlags().StringVarP(&format, "format", "f", "junit", "formating (default is junit)")
	return tmp
}
