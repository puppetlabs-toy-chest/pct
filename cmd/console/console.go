package console

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var ()

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "console [console_options]",
		Short: "(Experimental) Start a session of the puppet debugger console",
		Long: `
The pdk console runs an interactive session of the puppet debugger tool to test
out snippets of code, run language evaluations, datatype prototyping and much
more. A virtual playground for your puppet code!
For usage details see the puppet debugger docs at https://docs.puppet-debugger.com.`,
		RunE: utils.ExecutePDKCommand,
	}

	return tmp
}
