package test

import (
	"github.com/spf13/cobra"
)

var ()

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "test [flags]",
		Short: "Run tests",
		Long:  `Run tests`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return tmp
}
