package root

import (
	"os"

	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	LogLevel string
	debug    bool
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
		Use:   "pcv",
		Short: "pcv - Puppet Content Validator",
		Long:  `Puppet Content Validator (PCV) - create and execute a Validator container`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

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

	return tmp
}
