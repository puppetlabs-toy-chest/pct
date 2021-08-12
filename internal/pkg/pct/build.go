package pct

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/puppetlabs/pdkgo/internal/pkg/gzip"
	"github.com/puppetlabs/pdkgo/internal/pkg/tar"
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog/log"
)

var osUtil utils.OsUtil
var ioUtil utils.IoUtil
var tarUtil tar.TarHelpers
var gzipUtil gzip.GzipHelpers

func init() {
	osUtil = utils.OsUtilHelpersImpl{}
	ioUtil = utils.IoUtilHelpersImpl{}
	tarUtil = tar.TarHelpersImpl{}
	gzipUtil = gzip.GzipHelpersImpl{}
}

func Build(templatePath, targetDir string) (gzipArchiveFilePath string, err error) {
	// Check template dir exists
	if _, err := osUtil.Stat(templatePath); osUtil.IsNotExist(err) {
		return "", fmt.Errorf("No template directory at %v", templatePath)
	}

	// Check if pct-config.yml exists
	if _, err := osUtil.Stat(filepath.Join(templatePath, "pct-config.yml")); osUtil.IsNotExist(err) {
		return "", fmt.Errorf("No 'pct-config.yml' found in %v", templatePath)
	}

	// Check if content dir exists
	if _, err := osUtil.Stat(filepath.Join(templatePath, "content")); osUtil.IsNotExist(err) {
		return "", fmt.Errorf("No 'content' dir found in %v", templatePath)
	}

	// Create temp dir and TAR template there
	tempDir, err := ioUtil.TempDir()
	defer os.Remove(tempDir)

	if err != nil {
		log.Error().Msgf("Could not create tempdir to TAR template: %v", err)
		return "", err
	}

	tarArchiveFilePath, err := tarUtil.Tar(templatePath, tempDir)
	if err != nil {
		log.Error().Msgf("Could not TAR template (%v): %v", templatePath, err)
		return "", err
	}

	// GZIP the TAR created in the temp dir and output to the $MODULE_ROOT/pkg directory
	gzipArchiveFilePath, err = gzipUtil.Gzip(tarArchiveFilePath, filepath.Join(targetDir))
	if err != nil {
		log.Error().Msgf("Could not GZIP template TAR archive (%v): %v", tarArchiveFilePath, err)
		return "", err
	}

	return gzipArchiveFilePath, nil
}
