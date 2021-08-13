package new

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	localTemplateCache   string
	format               string
	selectedTemplate     string
	selectedTemplateInfo string
	listTemplates        bool
	targetName           string
	targetOutput         string
	pctApi               *pct.Pct
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

	// Configure PCT
	fs := afero.NewOsFs() // configure afero to use real filesystem
	pctApi = &pct.Pct{
		OsUtils: &utils.OsUtil{},
		Utils:   &utils.UtilsHelper{},
		AFS:     &afero.Afero{Fs: fs},
		IOFS:    &afero.IOFS{Fs: fs},
	}

	tmp.Flags().SortFlags = false

	tmp.Flags().StringVarP(&targetName, "name", "n", "", "the name for the created output.")
	tmp.Flags().StringVarP(&targetOutput, "output", "o", "", "location to place the generated output.")

	tmp.Flags().BoolVarP(&listTemplates, "list", "l", false, "list templates")
	err := tmp.RegisterFlagCompletionFunc("list", flagCompletion)
	cobra.CheckErr(err)

	tmp.Flags().StringVarP(&selectedTemplateInfo, "info", "i", "", "display the selected template's configuration and default values")
	err = tmp.RegisterFlagCompletionFunc("info", flagCompletion)
	cobra.CheckErr(err)

	tmp.Flags().StringVar(&format, "format", "table", "display output in table or json format")
	err = tmp.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var formats = []string{"table", "json"}
		return utils.Find(formats, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)

	tmp.Flags().StringVar(&localTemplateCache, "templatepath", "", "location of installed templates")
	err = viper.BindPFlag("templatepath", tmp.Flags().Lookup("templatepath"))
	cobra.CheckErr(err)

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
	// show available templates if user runs `pct new`
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
	if localTemplateCache == "" {
		err := preExecute(cmd, args)
		if err != nil {
			log.Error().Msgf("Unable to set template path: %s", err.Error())
			return nil, cobra.ShellCompDirectiveError
		}
	}
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	localTemplateCache = viper.GetString("templatepath")

	return completeName(localTemplateCache, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

func completeName(cache string, match string) []string {
	tmpls, _ := pctApi.List(cache, "")
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

	if listTemplates && selectedTemplateInfo == "" {
		tmpls, err := pctApi.List(localTemplateCache, selectedTemplate)
		if err != nil {
			return err
		}

		err = pctApi.FormatTemplates(tmpls, format)
		if err != nil {
			return err
		}

		return nil
	}

	if selectedTemplateInfo != "" {
		pctData, err := pctApi.GetInfo(localTemplateCache, selectedTemplateInfo)
		if err != nil {
			return err
		}
		log.Debug().Msgf("Template Defaults: %v", pctData.Defaults)
		defaultString := pctApi.DisplayDefaults(pctData.Defaults, format)
		fmt.Printf("%s\n", defaultString)

		return nil
	}

	_, err := pctApi.Get(localTemplateCache, selectedTemplate)
	if err != nil {
		return err
	}

	appVersionString := cmd.Parent().Version
	pdkInfo := getApplicationInfo(appVersionString)

	deployed := pctApi.Deploy(pct.DeployInfo{
		SelectedTemplate: selectedTemplate,
		TemplateCache:    localTemplateCache,
		TargetOutputDir:  targetOutput,
		TargetName:       targetName,
		PdkInfo:          pdkInfo,
	})

	err = pctApi.FormatDeployment(deployed, format)
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
