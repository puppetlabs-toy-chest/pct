package get

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PuppetVersion string

func CreateGetPuppetCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "puppet",
		Short: "Gets the Puppet version configured",
		Long: `Gets the Puppet version configured
		`,
		Run: func(cmd *cobra.Command, args []string) {
			LogPuppetVersion()
		},
	}

	return tmp
}

func LogPuppetVersion() {
	log.Info().Msgf("Puppet version is configured to: %s", viper.GetString(config.PuppetVersion))
}
