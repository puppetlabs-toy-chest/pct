package pct

import (
	"fmt"
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

type PctInstallerI interface {
	Install(templatePkg string, targetDir string) (string, error)
}

func (p *PctInstaller) Install(templatePkg string, targetDir string) (string, error) {

	if _, err := p.AFS.Stat(templatePkg); os.IsNotExist(err) {
		return "", fmt.Errorf("No template package at %v", templatePkg)
	}

	tempDir, err := p.AFS.TempDir("", "")
	defer func() {
		err := p.AFS.Remove(tempDir)
		log.Debug().Msgf("Failed to remove temp dir: %v", err)
	}()
	if err != nil {
		return "", fmt.Errorf("Could not create tempdir to gunzip template: %v", err)
	}

	err = p.Gunzip.Gunzip(templatePkg, tempDir)
	if err != nil {
		return "", fmt.Errorf("Could not extract TAR from GZIP (%v): %v", templatePkg, err)
	}

	tempFile := strings.TrimSuffix(filepath.Join(tempDir, filepath.Base(templatePkg)), `.gz`)
	if _, err = p.AFS.Stat(tempFile); os.IsNotExist(err) {
		return "", err
	}

	t, err := p.Tar.Untar(tempFile, targetDir)
	if err != nil {
		return "", fmt.Errorf("Could not UNTAR template (%v): %v", templatePkg, err)
	}

	return t, nil
}
