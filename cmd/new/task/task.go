package task

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug       bool
	description string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "task [options] <name>",
		Args:  cobra.MaximumNArgs(1),
		Short: "Create a new test for the object named <name>",
		Long:  `Create a new test for the object named <name>`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")

	tmp.Flags().StringVar(&description, "description", "", "A short description of the purpose of the task")

	return tmp
}
