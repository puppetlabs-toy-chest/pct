package api

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func GenerateCompletion(rootCmd *cobra.Command, shell string) {
	var err error

	switch shell {
	case "bash":
		err = rootCmd.GenBashCompletion(os.Stdout)
	case "fish":
		err = rootCmd.GenFishCompletion(os.Stdout, true)
	case "pwsh":
		err = rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
	case "zsh":
		err = rootCmd.GenZshCompletion(os.Stdout)
	default:
		log.Printf("unsupported shell type %q", shell)
	}

	if err != nil {
		log.Fatal(err)
	}
}
