package bundle

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var ()

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "bundle [bundler_options]",
		Short: "(Experimental) Command pass-through to bundler",
		Long: `[experimental] For advanced users, pdk bundle runs arbitrary commands in the bundler environment that pdk manages.
		Careless use of this command can lead to errors that pdk can't help recover from.`,
		RunE: utils.ExecutePDKCommand,
	}

	return tmp
}
