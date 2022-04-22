package completion

import (
	"log"
	"os"

	"github.com/puppetlabs/pct/pkg/telemetry"
	"github.com/spf13/cobra"
)

func CreateCompletionCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completions for the chosen shell",
		Long: `To load completions:

Bash:

	$ source <(pct completion bash)

	# To load completions for each session, execute once:
	# Linux:
	$ pct completion bash > /etc/bash_completion.d/pct
	# macOS:
	$ pct completion bash > /usr/local/etc/bash_completion.d/pct

Zsh:

	# If shell completion is not already enabled in your environment,
	# you will need to enable it.  You can execute the following once:

	$ echo "autoload -U compinit; compinit" >> ~/.zshrc

	# To load completions for each session, execute once:
	$ pct completion zsh > "${fpath[1]}/_pct"

	# You will need to start a new shell for this setup to take effect.

fish:

	$ pct completion fish | source

	# To load completions for each session, execute once:
	$ pct completion fish > ~/.config/fish/completions/pct.fish

PowerShell:

	PS> pct completion powershell | Out-String | Invoke-Expression

	# To load completions for every new session, run:
	PS> pct completion powershell > pct.ps1
	# and source this file from your PowerShell profile.`,
		ValidArgs: []string{"bash", "fish", "pwsh", "zsh"},
		Args:      cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_, span := telemetry.NewSpan(cmd.Context(), "completion")
			defer telemetry.EndSpan(span)
			telemetry.AddStringSpanAttribute(span, "name", "completion")

			var err error
			switch args[0] {
			case "bash":
				err = cmd.Root().GenBashCompletion(os.Stdout)
			case "fish":
				err = cmd.Root().GenFishCompletion(os.Stdout, true)
			case "pwsh":
				err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			case "zsh":
				err = cmd.Root().GenZshCompletion(os.Stdout)
			default:
				log.Printf("unsupported shell type %q", args[0])
			}

			if err != nil {
				log.Fatal(err)
			}
		},
	}
	return tmp
}
