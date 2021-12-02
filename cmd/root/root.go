package root

import (
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/puppetlabs/pdkgo/pkg/utils"
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
		Use:   "pct",
		Short: "pct - Puppet Content Templater",
		Long:  `Puppet Content Templater (PCT) - Create a range of Puppet content from templates`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	tmp.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pdk.yaml)")

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

// Returns the cobra command called, e.g. new or install
// and also the fully formatted command as passed with arguments/flags.
// Idea borrowed from carolynvs/porter:
// https://github.com/carolynvs/porter/blob/ccca10a63627e328616c1006600153da8411a438/cmd/porter/main.go
func GetCalledCommand(cmd *cobra.Command) (string, string) {
	if len(os.Args) < 2 {
		return "", ""
	}

	calledCommandStr := os.Args[1]

	// Also figure out the full called command from the CLI
	// Is there anything sensitive you could leak here? 🤔
	calledCommandArgs := strings.Join(os.Args[1:], " ")

	return calledCommandStr, calledCommandArgs
}
