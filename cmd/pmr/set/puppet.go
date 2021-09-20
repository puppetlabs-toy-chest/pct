package set

import (
	"github.com/puppetlabs/pdkgo/cmd/pmr/get"
	"github.com/puppetlabs/pdkgo/internal/pkg/config"
	"github.com/puppetlabs/pdkgo/internal/pkg/puppet"
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PuppetVersion string

func CreateSetPuppetCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "puppet",
		Short: "Sets the Puppet version to use with the PDK",
		Long: `Sets the Puppet version to use with the PDK.
		REQUIRES DOCKER.
		This will:
			* Download a docker image containing the PDK toolchain for the version of Puppet you set.
			* Stand up the container
			* Link the PDK to the container for all subsequent functions.
			* The setting will persist until called again.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			configurePuppet(args)
		},
		ValidArgs: []string{"5.5.0", "6.15.0", "7.5.0"},
		Args:      cobra.ExactValidArgs(1),
	}

	tmp.PersistentFlags().StringVar(&PuppetVersion, "version", zerolog.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
	tmp.RegisterFlagCompletionFunc("version", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var levels = []string{"debug", "info", "warn", "error", "fatal", "panic"}
		return utils.Find(levels, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})

	return tmp
}

func configurePuppet(args []string) {
	viper.Set(config.PuppetVersion, args[0])
	if err := viper.WriteConfig(); err != nil {
		log.Error().Msg(err.Error())
	}

	puppet.StopContainer()
	puppet.StartContainer(viper.GetString(config.PuppetVersion))
	get.LogPuppetVersion()
}
