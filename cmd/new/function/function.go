package function

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug        bool
	functionType string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "function [options] <name>",
		Args:  cobra.MaximumNArgs(1),
		Short: "Create a new custom function named <name> using given options",
		Long:  `Create a new custom function named <name> using given options`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	tmp.Flags().StringVar(&functionType, "type", "", "The function type, (native or v4) (default: native)")
	tmp.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var levels = []string{"native", "v4"}
		return utils.Find(levels, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})

	return tmp
}
