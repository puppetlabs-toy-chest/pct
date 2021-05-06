package new

import (
	"path/filepath"
	"strings"

	"github.com/puppetlabs/pdkgo/internal/pkg/pct"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	localTemplateCache string
	jsonOutput         bool
	selectedTemplate   string
	listTemplates      bool
	targetName         string
	targetOutput       string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:               "new <template> [args] [flags]",
		Short:             "Creates a Puppet project or other artifact based on a template",
		Long:              `Creates a Puppet project or other artifact based on a template`,
		Args:              validateArgCount,
		ValidArgsFunction: flagCompletion,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			log.Trace().Msg("PreRunE")
			localTemplateCache = viper.GetString("templatepath")
			return nil
		},
		RunE: execute,
	}

	tmp.Flags().StringVar(&localTemplateCache, "templatepath", "", "Location of installed templates")
	viper.BindPFlag("templatepath", tmp.Flags().Lookup("templatepath")) //nolint:errcheck
	home, _ := homedir.Dir()
	viper.SetDefault("templatepath", filepath.Join(home, ".pdk", "pct"))

	tmp.Flags().BoolVarP(&listTemplates, "list", "l", false, "list templates")
	tmp.RegisterFlagCompletionFunc("list", flagCompletion) //nolint:errcheck

	tmp.Flags().StringVarP(&targetName, "name", "n", "", "the name for the created output. (default is the name of the current directory)")
	tmp.Flags().StringVarP(&targetOutput, "output", "o", "", "location to place the generated output. (default is the current directory)")
	tmp.Flags().BoolVar(&jsonOutput, "json", false, "json output")
	return tmp
}

func validateArgCount(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && !listTemplates {
		listTemplates = true
	}

	if targetName == "" && len(args) == 2 {
		targetName = args[1]
	}

	if len(args) >= 1 {
		selectedTemplate = args[0]
	}

	return nil
}

func flagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	localTemplateCache = viper.GetString("templatepath")
	return completeName(localTemplateCache, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

func completeName(cache string, match string) []string {
	tmpls, _ := pct.List(cache, "")
	var names []string
	for _, tmpl := range tmpls {
		if strings.HasPrefix(tmpl.Name, match) {
			m := tmpl.Name + "\t" + tmpl.Display
			names = append(names, m)
		}
	}
	return names
}

func getApplicationInfo(appVersionString string) pct.PDKInfo {
	info := strings.Split(appVersionString, "\n")[0]
	appInfo := strings.Split(info, " ")
	ver := appInfo[1]
	commit := appInfo[2]
	buildDate := appInfo[3]

	pdkInfo := pct.PDKInfo{
		Version:   ver,
		Commit:    commit,
		BuildDate: buildDate,
	}
	return pdkInfo
}

func execute(cmd *cobra.Command, args []string) error {
	log.Trace().Msg("Run")
	log.Trace().Msgf("Template path: %v", localTemplateCache)
	log.Trace().Msgf("Selected template: %v", selectedTemplate)

	if listTemplates {
		tmpls, err := pct.List(localTemplateCache, selectedTemplate)
		if err != nil {
			return err
		}

		err = pct.FormatTemplates(tmpls, jsonOutput)
		if err != nil {
			return err
		}

		return nil
	}

	appVersionString := cmd.Parent().Version
	pdkInfo := getApplicationInfo(appVersionString)

	deployed := pct.Deploy(
		selectedTemplate,
		localTemplateCache,
		targetOutput,
		targetName,
		pdkInfo,
	)

	err := pct.FormatDeployment(deployed, jsonOutput)
	if err != nil {
		return err
	}

	return nil
}
