package get

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CreateGetCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "get",
		Short: "Gets PDK configuration",
		Long:  `Gets PDK configuration`,
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

	tmp.AddCommand(CreateGetPuppetCommand())

	return tmp
}
