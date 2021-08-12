package pct

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/puppetlabs/pdkgo/internal/pkg/gzip"
	"github.com/puppetlabs/pdkgo/internal/pkg/tar"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type PctInstaller struct {
	Tar    tar.TarI
	Gunzip gzip.GunzipI
	AFS    *afero.Afero
	IOFS   *afero.IOFS
}

func (p *PctInstaller) Install(tarFile string, targetDir string) (string, error) {
	tempDir, err := p.AFS.TempDir("", "")
	defer func() {
		err := p.AFS.Remove(tempDir)
		log.Debug().Msgf("Failed to remove temp dir: %v", err)
	}()
	if err != nil {
		log.Error().Msgf("Could not create tempdir to gunzip template: %v", err)
		return "", err
	}

	err = p.Gunzip.Gunzip(tarFile, tempDir)
	if err != nil {
		log.Error().Msgf("Could not GZIP template TAR archive (%v): %v", tempDir, err)
		return "", err
	}

	tempFile := strings.TrimSuffix(filepath.Join(tempDir, filepath.Base(tarFile)), `.gz`)
	if _, err = p.AFS.Stat(tempFile); os.IsNotExist(err) {
		return "", err
	}

	t, err := p.Tar.Untar(tempFile, targetDir)
	if err != nil {
		log.Error().Msgf("Could not UNTAR template (%v): %v", tarFile, err)
		return "", err
	}

	return t, nil
}
