package install

import (
	"fmt"

	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/puppetlabs/pdkgo/pkg/telemetry"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type InstallCommand struct {
	TemplatePkgPath string
	InstallPath     string
	Force           bool
	PctInstaller    pct.PctInstallerI
}

type InstallCommandI interface {
	CreateCommand() *cobra.Command
}

func (ic *InstallCommand) CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:     "install [flags]",
		Short:   "Installs a template package (in tar.gz format)",
		Long:    `Installs a template package (in tar.gz format) to the default or specified template path`,
		PreRunE: ic.preExecute,
		RunE:    ic.executeInstall,
	}
	tmp.Flags().StringVar(&ic.InstallPath, "templatepath", "", "location of installed templates")
	err := viper.BindPFlag("templatepath", tmp.Flags().Lookup("templatepath"))
	tmp.Flags().BoolVarP(&ic.Force, "force", "f", false, "Forces the install of a template without error, if it already exists. ")

	cobra.CheckErr(err)

	return tmp
}

func (ic *InstallCommand) executeInstall(cmd *cobra.Command, args []string) error {
	_, span := telemetry.NewSpan(cmd.Context(), "install")
	defer telemetry.EndSpan(span)
	telemetry.AddStringSpanAttribute(span, "name", "install")

	templateInstallationPath, err := ic.PctInstaller.Install(ic.TemplatePkgPath, ic.InstallPath, ic.Force)
	if err != nil {
		telemetry.RecordSpanError(span, err)
		return err
	}
	log.Info().Msgf("Template installed to %v", templateInstallationPath)
	return nil
}

func (ic *InstallCommand) setInstallPath() error {
	if ic.InstallPath == "" {
		ic.InstallPath = viper.GetString("templatepath")
		if ic.InstallPath == "" {
			defaultTemplatePath, err := utils.GetDefaultTemplatePath()
			if err != nil {
				return fmt.Errorf("Could not determine location to install template: %v", err)
			}
			ic.InstallPath = defaultTemplatePath
		}
	}
	return nil
}

func (ic *InstallCommand) preExecute(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Path to template package (tar.gz) should be first argument")
	}

	if len(args) == 1 {
		ic.TemplatePkgPath = args[0]
		return ic.setInstallPath()
	}

	if len(args) > 1 {
		return fmt.Errorf("Incorrect number of arguments; path to template package (tar.gz) should be first argument")
	}

	return nil
}
