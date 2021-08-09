package build

import (
	"os"
	"path/filepath"

	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	sourceDir string
	targetDir string
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:     "build [flags]",
		Short:   "Builds a package from the template",
		Long:    `Builds a package from the template. Assumes the current working directory is the template you wish to package`,
		PreRunE: preExecute,
		RunE:    execute,
	}

	tmp.Flags().StringVar(&sourceDir, "sourcedir", "", "The template directory you wish to package up")
	tmp.Flags().StringVar(&targetDir, "targetdir", "", "The target directory where you want the packaged template to be output to")

	return tmp
}

func preExecute(cmd *cobra.Command, args []string) error {

	wd, err := os.Getwd()
	log.Info().Msgf("WD: %v", wd)

	if (sourceDir == "" || targetDir == "") && err != nil {
		return err
	}

	if sourceDir == "" {
		sourceDir = wd
	}

	if targetDir == "" {
		targetDir = filepath.Join(wd, "pkg")
	}

	return nil
}

func execute(cmd *cobra.Command, args []string) error {
	gzipArchiveFilePath, err := pct.Build(sourceDir, targetDir)
	log.Info().Msgf("Template output to %v", gzipArchiveFilePath)
	return err
}
