package new

import (
	"path/filepath"
	"strings"

	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	localTemplateCache string
	format             string
	selectedTemplate   string
	listTemplates      bool
	targetName         string
	targetOutput       string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:               "new <template> [flags]",
		Short:             "Creates a Puppet project or other artifact based on a template",
		Long:              `Creates a Puppet project or other artifact based on a template`,
		Args:              validateArgCount,
		ValidArgsFunction: flagCompletion,
		PreRunE:           preExecute,
		RunE:              execute,
	}

	tmp.Flags().SortFlags = false

	tmp.Flags().StringVarP(&targetName, "name", "n", "", "the name for the created output.")
	tmp.Flags().StringVarP(&targetOutput, "output", "o", "", "location to place the generated output.")

	tmp.Flags().BoolVarP(&listTemplates, "list", "l", false, "list templates")
	tmp.RegisterFlagCompletionFunc("list", flagCompletion) //nolint:errcheck

	tmp.Flags().StringVar(&format, "format", "table", "display output in table or json format")
	tmp.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var formats = []string{"table", "json"}
		return utils.Find(formats, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})

	tmp.Flags().StringVar(&localTemplateCache, "templatepath", "", "location of installed templates")
	viper.BindPFlag("templatepath", tmp.Flags().Lookup("templatepath")) //nolint:errcheck
	return tmp
}

func preExecute(cmd *cobra.Command, args []string) error {
	defaultTemplatePath, err := getDefaultTemplatePath()
	if err != nil {
		return err
	}

	viper.SetDefault("templatepath", defaultTemplatePath)
	localTemplateCache = viper.GetString("templatepath")
	return nil
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
		if strings.HasPrefix(tmpl.Id, match) {
			m := tmpl.Id + "\t" + tmpl.Display
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

		err = pct.FormatTemplates(tmpls, format)
		if err != nil {
			return err
		}

		return nil
	}

	_, err := pct.Get(localTemplateCache, selectedTemplate)
	if err != nil {
		return err
	}

	appVersionString := cmd.Parent().Version
	pdkInfo := getApplicationInfo(appVersionString)

	deployed := pct.Deploy(pct.DeployInfo{
		SelectedTemplate: selectedTemplate,
		TemplateCache:    localTemplateCache,
		TargetOutputDir:  targetOutput,
		TargetName:       targetName,
		PdkInfo:          pdkInfo,
	})

	err = pct.FormatDeployment(deployed, format)
	if err != nil {
		return err
	}

	return nil
}

func getDefaultTemplatePath() (string, error) {
	execDir, err := os.Executable()
	if err != nil {
		return "", err
	}

	defaultTemplatePath := filepath.Join(filepath.Dir(execDir), "templates")
	log.Trace().Msgf("Default template path: %v", defaultTemplatePath)
	return defaultTemplatePath, nil
}
