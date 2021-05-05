package module

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug         bool
	fullInterview bool
	skipInterview bool
	license       string
	templateRef   string
	templateURL   string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "module [options] [module_name] [target_dir]",
		Args:  cobra.MaximumNArgs(2),
		Short: "Create a new module named [module_name] using given options",
		Long:  `Create a new module named [module_name] using given options`,
		RunE:  utils.ExecutePDKCommand,
	}
	tmp.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	tmp.Flags().BoolVar(&fullInterview, "full-interview", false, "Enable debug output")
	tmp.Flags().BoolVar(&skipInterview, "skip-interview", false, "Enable debug output")
	tmp.Flags().StringVar(&license, "license", "", "The function type, (native or v4) (default: native)")
	tmp.Flags().StringVar(&templateRef, "template-ref", "", "The function type, (native or v4) (default: native)")
	tmp.Flags().StringVar(&templateURL, "template-url", "", "The function type, (native or v4) (default: native)")

	return tmp
}
