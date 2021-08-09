package pct

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/puppetlabs/pdkgo/internal/pkg/gzip"
	"github.com/puppetlabs/pdkgo/internal/pkg/tar"
	"github.com/rs/zerolog/log"
)

var untarUtil tar.UntarHelpers
var gunzipUtil gzip.GunzipHelpers

func init() {
	untarUtil = tar.UntarHelpersImpl{}
	gunzipUtil = gzip.GunzipHelpersImpl{}
}

func Install(tarFile string, targetDir string) (string, error) {
	tempDir, err := ioUtil.TempDir()
	defer os.Remove(tempDir)
	if err != nil {
		log.Error().Msgf("Could not create tempdir to gunzip template: %v", err)
		return "", err
	}

	err = gunzipUtil.Gunzip(tarFile, tempDir)
	if err != nil {
		log.Error().Msgf("Could not GZIP template TAR archive (%v): %v", tempDir, err)
		return "", err
	}

	tempFile := strings.TrimSuffix(filepath.Join(tempDir, filepath.Base(tarFile)), `.gz`)
	if _, err := osUtil.Stat(tempFile); osUtil.IsNotExist(err) {
		return "", err
	}

	t, err := untarUtil.Untar(tempFile, targetDir)
	if err != nil {
		log.Error().Msgf("Could not UNTAR template (%v): %v", tarFile, err)
		return "", err
	}

	return t, nil
}
