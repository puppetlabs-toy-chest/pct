package config

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	debug     bool
	force     bool
	valueType string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "config [subcommand] [options]",
		Short: "Set or update the configuration for <name>",
		Long:  `Set or update the configuration for <name>`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	tmp.Flags().BoolVar(&force, "force", false, "Force the configuration setting to be overwitten")

	tmp.Flags().StringVar(&valueType, "type", "", "The type of value to set. Acceptable values: 'array', 'boolean', 'number', 'string'")
	tmp.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		switch name {
		case "as":
			name = "type"
		}
		return pflag.NormalizedName(name)
	})
	tmp.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var levels = []string{"array", "boolean", "number", "string"}
		return utils.Find(levels, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
	return tmp
}
