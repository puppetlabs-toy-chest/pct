package root

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	tmp.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pdk.yaml)")

	tmp.PersistentFlags().StringVar(&LogLevel, "log-level", zerolog.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
	tmp.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var levels = []string{"debug", "info", "warn", "error", "fatal", "panic"}
		return find(levels, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})

	return tmp
}

func InitConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, _ := homedir.Dir()
		viper.SetConfigName(".pdk")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, ".config"))
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Trace().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}
}

// finds a string present in a slice
func find(s []string, str string) []string {
	var matches []string
	if contains(s, str) {
		matches = append(matches, str)
	}
	return matches
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
