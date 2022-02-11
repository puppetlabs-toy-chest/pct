package pct

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/puppetlabs/pdkgo/pkg/config_processor"
	"github.com/puppetlabs/pdkgo/pkg/gzip"
	"github.com/puppetlabs/pdkgo/pkg/tar"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type BuilderI interface {
	Build(templatePath, targetDir string) (gzipArchiveFilePath string, err error)
}

type Builder struct {
	Tar             tar.TarI
	Gzip            gzip.GzipI
	AFS             *afero.Afero
	ConfigProcessor config_processor.ConfigProcessorI
}

func (b *Builder) Build(templatePath, targetDir string) (gzipArchiveFilePath string, err error) {
	// Check template dir exists
	if _, err := b.AFS.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("No template directory at %v", templatePath)
	}

	// Check if pct-config.yml exists
	if _, err := b.AFS.Stat(filepath.Join(templatePath, "pct-config.yml")); os.IsNotExist(err) {
		return "", fmt.Errorf("No 'pct-config.yml' found in %v", templatePath)
	}

	err = b.ConfigProcessor.CheckConfig(filepath.Join(templatePath, "pct-config.yml"))
	if err != nil {
		return "", fmt.Errorf("Invalid config: %v", err.Error())
	}

	// Check if content dir exists
	if _, err := b.AFS.Stat(filepath.Join(templatePath, "content")); os.IsNotExist(err) {
		return "", fmt.Errorf("No 'content' dir found in %v", templatePath)
	}

	// Create temp dir and TAR template there
	tempDir, err := b.AFS.TempDir("", "")
	defer os.Remove(tempDir)

	if err != nil {
		log.Error().Msgf("Could not create tempdir to TAR template: %v", err)
		return "", err
	}

	tarArchiveFilePath, err := b.Tar.Tar(templatePath, tempDir)
	if err != nil {
		log.Error().Msgf("Could not TAR template (%v): %v", templatePath, err)
		return "", err
	}

	// GZIP the TAR created in the temp dir and output to the $MODULE_ROOT/pkg directory
	gzipArchiveFilePath, err = b.Gzip.Gzip(tarArchiveFilePath, targetDir)
	if err != nil {
		log.Error().Msgf("Could not GZIP template TAR archive (%v): %v", tarArchiveFilePath, err)
		return "", err
	}

	return gzipArchiveFilePath, nil
}
