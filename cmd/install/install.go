package install

import (
	"fmt"

	"github.com/puppetlabs/pdkgo/internal/pkg/gzip"
	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/puppetlabs/pdkgo/internal/pkg/tar"
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	templatePkgPath string
	installPath     string
	installer       *pct.PctInstaller
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:     "install [flags]",
		Short:   "Installs a template package (in tar.gz) format",
		Long:    `Installs a template package (in tar.gz) format to the default or specified template path`,
		PreRunE: preExecute,
		RunE:    executeInstall,
	}

	tmp.Flags().StringVar(&installPath, "templatepath", "", "location of installed templates")
	err := viper.BindPFlag("templatepath", tmp.Flags().Lookup("templatepath"))
	cobra.CheckErr(err)

	fs := afero.NewOsFs() // configure afero to use real filesystem
	installer = &pct.PctInstaller{
		Tar:    &tar.Tar{AFS: &afero.Afero{Fs: fs}},
		Gunzip: &gzip.Gunzip{AFS: &afero.Afero{Fs: fs}},
		AFS:    &afero.Afero{Fs: fs},
		IOFS:   &afero.IOFS{Fs: fs},
	}

	return tmp
}

func preExecute(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Must specify a path to a template package")
	}

	if len(args) == 1 {
		templatePkgPath = args[0]
		if installPath == "" {
			installPath = viper.GetString("templatepath")
			if installPath == "" {
				defaultTemplatePath, err := utils.GetDefaultTemplatePath()
				if err != nil {
					return fmt.Errorf("Could not determine location to install template: %v", err)
				}
				installPath = defaultTemplatePath
			}
		}
	}

	return nil
}

func executeInstall(cmd *cobra.Command, args []string) error {
	templateInstallPath, err := installer.Install(templatePkgPath, installPath)

	if err != nil {
		return err
	}

	log.Info().Msgf("Template installed to %v", templateInstallPath)
	return nil
}
