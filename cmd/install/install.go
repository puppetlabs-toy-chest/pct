package install

import (
	"fmt"

	"github.com/spf13/afero"

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
	GitUri          string
	AFS             *afero.Afero
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
	tmp.Flags().StringVar(&ic.GitUri, "git-uri", "", "Installs a template package from a remote git repository.")

	cobra.CheckErr(err)

	return tmp
}

func (ic *InstallCommand) executeInstall(cmd *cobra.Command, args []string) error {
	_, span := telemetry.NewSpan(cmd.Context(), "install")
	defer telemetry.EndSpan(span)
	telemetry.AddStringSpanAttribute(span, "name", "install")

	templateInstallationPath := ""
	var err error = nil
	if ic.GitUri != "" { // For cloning a template
		// Create temp folder
		tempDir, dirErr := ic.AFS.TempDir("", "")
		defer func() {
			dirErr := ic.AFS.Remove(tempDir)
			if dirErr != nil {
				log.Error().Msgf("Failed to remove temp dir: %v", dirErr)
			}
		}()
		if dirErr != nil {
			return fmt.Errorf("Could not create tempdir to clone template to: %v", err)
		}
		templateInstallationPath, err = ic.PctInstaller.InstallClone(ic.GitUri, ic.InstallPath, tempDir, ic.Force)
	} else { // For downloading and/or locally installing a template
		templateInstallationPath, err = ic.PctInstaller.Install(ic.TemplatePkgPath, ic.InstallPath, ic.Force)
	}

	if err != nil {
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
		if ic.GitUri != "" {
			return ic.setInstallPath()
		}
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
