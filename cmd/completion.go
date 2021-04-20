package cmd

import (
	"fmt"

	"github.com/puppetlabs/pdkgo/pkg/api"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completions for the chosen shell",
	Long: `Generate shell completions for the chosen shell.

	Example usage:
		pdkgo completion bash
		or
		pdkgo completion pwsh`,
	ValidArgs: []string{"bash", "fish", "pwsh", "zsh"},
	Args:      cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api.GenerateCompletion(rootCmd, args[0])
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)

	shellMsg := fmt.Sprintf("shell to generate script for: %v", completionCmd.ValidArgs)
	completionCmd.Flags().StringP("shell", "s", "", shellMsg)
}
