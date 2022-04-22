package build

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/puppetlabs/pct/pkg/config_processor"
	"github.com/puppetlabs/pct/pkg/gzip"
	"github.com/puppetlabs/pct/pkg/tar"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type BuilderI interface {
	Build(sourceDir, targetDir string) (gzipArchiveFilePath string, err error)
}

type Builder struct {
	Tar             tar.TarI
	Gzip            gzip.GzipI
	AFS             *afero.Afero
	ConfigProcessor config_processor.ConfigProcessorI
	ConfigFile      string
}

func (b *Builder) Build(sourceDir, targetDir string) (gzipArchiveFilePath string, err error) {
	// Check project dir exists
	if _, err := b.AFS.Stat(sourceDir); os.IsNotExist(err) {
		return "", fmt.Errorf("No project directory at %v", sourceDir)
	}

	// Check if config file exists
	if _, err := b.AFS.Stat(filepath.Join(sourceDir, b.ConfigFile)); os.IsNotExist(err) {
		return "", fmt.Errorf("No '%v' found in %v", b.ConfigFile, sourceDir)
	}

	err = b.ConfigProcessor.CheckConfig(filepath.Join(sourceDir, b.ConfigFile))
	if err != nil {
		return "", fmt.Errorf("Invalid config: %v", err.Error())
	}

	// Check if content dir exists
	if _, err := b.AFS.Stat(filepath.Join(sourceDir, "content")); os.IsNotExist(err) {
		return "", fmt.Errorf("No 'content' dir found in %v", sourceDir)
	}

	// Create temp dir and TAR project there
	tempDir, err := b.AFS.TempDir("", "")
	defer os.Remove(tempDir)

	if err != nil {
		log.Error().Msgf("Could not create tempdir to TAR project: %v", err)
		return "", err
	}

	tarArchiveFilePath, err := b.Tar.Tar(sourceDir, tempDir)
	if err != nil {
		log.Error().Msgf("Could not TAR project (%v): %v", sourceDir, err)
		return "", err
	}

	// GZIP the TAR created in the temp dir and output to the /pkg directory in the target directory
	gzipArchiveFilePath, err = b.Gzip.Gzip(tarArchiveFilePath, targetDir)
	if err != nil {
		log.Error().Msgf("Could not GZIP project TAR archive (%v): %v", tarArchiveFilePath, err)
		return "", err
	}

	return gzipArchiveFilePath, nil
}
