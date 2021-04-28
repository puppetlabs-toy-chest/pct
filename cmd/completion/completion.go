package completion

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func CreateCompletionCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:       "completion",
		Short:     "Generate shell completions for the chosen shell",
		Long:      `To load completions:

Bash:

	$ source <(pdk completion bash)

	# To load completions for each session, execute once:
	# Linux:
	$ pdk completion bash > /etc/bash_completion.d/pdk
	# macOS:
	$ pdk completion bash > /usr/local/etc/bash_completion.d/pdk

Zsh:

	# If shell completion is not already enabled in your environment,
	# you will need to enable it.  You can execute the following once:

	$ echo "autoload -U compinit; compinit" >> ~/.zshrc

	# To load completions for each session, execute once:
	$ pdk completion zsh > "${fpath[1]}/_pdk"

	# You will need to start a new shell for this setup to take effect.

fish:

	$ pdk completion fish | source

	# To load completions for each session, execute once:
	$ pdk completion fish > ~/.config/fish/completions/pdk.fish

PowerShell:

	PS> pdk completion powershell | Out-String | Invoke-Expression

	# To load completions for every new session, run:
	PS> pdk completion powershell > pdk.ps1
	# and source this file from your PowerShell profile.`,
		ValidArgs: []string{"bash", "fish", "pwsh", "zsh"},
		Args:      cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
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
