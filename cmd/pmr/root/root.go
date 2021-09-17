package root

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile            string
	LogLevel           string
	LocalTemplateCache string

	debug bool
	// format string
)

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	lvl, err := zerolog.ParseLevel(LogLevel)
	if err != nil {
		panic("Could not initialize zerolog")
	}

	zerolog.SetGlobalLevel(lvl)

	if lvl == zerolog.InfoLevel {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	}

	log.Trace().Msg("Initialized zerolog")
}

func CreateRootCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "prm",
		Short: "prm - Puppet Runtime Manager",
		Long:  `Puppet Runtime Manager (PMR) - install and swap between Puppet Runtimes`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	tmp.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pmr.yaml)")

	tmp.PersistentFlags().StringVar(&LogLevel, "log-level", zerolog.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
	err := tmp.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var levels = []string{"debug", "info", "warn", "error", "fatal", "panic"}
		return utils.Find(levels, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)

	tmp.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug output")
	// tmp.PersistentFlags().StringVarP(&format, "format", "f", "junit", "formating (default is junit)")

	return tmp
}

func InitConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, _ := homedir.Dir()
		viper.SetConfigName(".pmr")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, ".config"))
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Trace().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}
}
