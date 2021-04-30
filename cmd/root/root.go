package root

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/puppetlabs/pdkgo/cmd/get"
	"github.com/puppetlabs/pdkgo/internal/pkg/config"
)

var (
	cfgFile            string
	LogLevel           string
	LocalTemplateCache string
)

func CreateRootCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "pdk",
		Short: "pdk - Puppet Development Kit",
		Long:  `Puppet development tooling, content creation, and testing framework`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

			lvl, err := zerolog.ParseLevel(LogLevel)
			if err != nil {
				return err
			}

			zerolog.SetGlobalLevel(lvl)

			log.Logger = log.
				Output(zerolog.ConsoleWriter{Out: os.Stdout}).
				With().Caller().Logger()
			log.Trace().Msg("PersistentPreRunE")

			return nil
		},
	}
	tmp.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pdk.yaml)")

	tmp.PersistentFlags().StringVar(&LogLevel, "log-level", zerolog.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
	tmp.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var levels = []string{"debug", "info", "warn", "error", "fatal", "panic"}
		return utils.Find(levels, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})

	return tmp
}

func InitConfig() {
	viper.SetDefault(config.PuppetVersion, "7")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, _ := homedir.Dir()
		viper.SetConfigName(".pdk.yaml")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		defaultConfigFile := fmt.Sprintf("%s/.pdk.yaml", home)
		if err:= viper.SafeWriteConfigAs(defaultConfigFile); err != nil {
			log.Error().Msgf("Failed to create config at `%s`: %s", defaultConfigFile,  err.Error())
		}
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Trace().Msgf("Using config file: %s", viper.ConfigFileUsed())
		get.LogPuppetVersion()
	}
}
