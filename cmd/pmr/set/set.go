package set

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CreateSetCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "set",
		Short: "Sets configuration for PDK",
		Long:  `Sets configuration for PDK`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("This command does not take arguments, got %d", len(args))
			}

			return nil
		},
	}

	tmp.AddCommand(CreateSetPuppetCommand())

	return tmp
}
