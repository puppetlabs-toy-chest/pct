package build

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/puppetlabs/pct/pkg/build"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type BuildCommand struct {
	SourceDir   string
	TargetDir   string
	ProjectType string
	Builder     build.BuilderI
}

type BuildCommandI interface {
	CreateCommand() *cobra.Command
}

func (bc *BuildCommand) CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:     "build [flags]",
		Short:   fmt.Sprintf("Builds a package from the %s project", bc.ProjectType),
		Long:    fmt.Sprintf("Builds a package from the %s project. Assumes the current working directory is the template you wish to package", bc.ProjectType),
		PreRunE: bc.preExecute,
		RunE:    bc.execute,
	}

	tmp.Flags().StringVar(&bc.SourceDir, "sourcedir", "", fmt.Sprintf("The %s project directory you wish to package up", bc.ProjectType))
	tmp.Flags().StringVar(&bc.TargetDir, "targetdir", "", fmt.Sprintf("The target directory where you want the packaged %s project to be output to", bc.ProjectType))

	return tmp
}

func (bc *BuildCommand) preExecute(cmd *cobra.Command, args []string) error {

	wd, err := os.Getwd()
	log.Trace().Msgf("WD: %v", wd)

	if (bc.SourceDir == "" || bc.TargetDir == "") && err != nil {
		return err
	}

	if bc.SourceDir == "" {
		bc.SourceDir = wd
	}

	bc.SourceDir = filepath.Clean(bc.SourceDir)

	if bc.TargetDir == "" {
		bc.TargetDir = filepath.Join(wd, "pkg")
	}

	return nil
}

func (bc *BuildCommand) execute(cmd *cobra.Command, args []string) error {
	gzipArchiveFilePath, err := bc.Builder.Build(bc.SourceDir, bc.TargetDir)

	if err != nil {
		return fmt.Errorf("`sourcedir` is not a valid %s project: %s", bc.ProjectType, err.Error())
	}
	log.Info().Msgf("Packaged %s output to %v", bc.ProjectType, gzipArchiveFilePath)
	return nil
}
