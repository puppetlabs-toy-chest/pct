package pct

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/puppetlabs/pdkgo/internal/pkg/gzip"
	"github.com/puppetlabs/pdkgo/internal/pkg/tar"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type PctInstaller struct {
	Tar    tar.TarI
	Gunzip gzip.GunzipI
	AFS    *afero.Afero
	IOFS   *afero.IOFS
}

type PctInstallerI interface {
	Install(templatePkg string, targetDir string, force bool) (string, error)
}

func (p *PctInstaller) Install(templatePkg string, targetDir string, force bool) (string, error) {

	if _, err := p.AFS.Stat(templatePkg); os.IsNotExist(err) {
		return "", fmt.Errorf("No template package at %v", templatePkg)
	}

	// create a temporary Directory to extract the tar.gz to
	tempDir, err := p.AFS.TempDir("", "")
	defer func() {
		err := p.AFS.Remove(tempDir)
		log.Debug().Msgf("Failed to remove temp dir: %v", err)
	}()
	if err != nil {
		return "", fmt.Errorf("Could not create tempdir to gunzip template: %v", err)
	}

	// gunzip the tar.gz to created tempdir
	tarfile, err := p.Gunzip.Gunzip(templatePkg, tempDir)
	if err != nil {
		return "", fmt.Errorf("Could not extract TAR from GZIP (%v): %v", templatePkg, err)
	}

	// untar the above archive to the temp dir
	untarPath, err := p.Tar.Untar(tarfile, tempDir)
	if err != nil {
		return "", fmt.Errorf("Could not UNTAR template (%v): %v", templatePkg, err)
	}

	// determine the properties of the template
	info, err := p.readConfig(filepath.Join(untarPath, "pct-config.yml"))
	if err != nil {
		return "", fmt.Errorf("Invalid config: %v", err.Error())
	}

	namespacedPath, err := p.setupTemplateNamespace(targetDir, info, untarPath, force)
	if err != nil {
		return "", fmt.Errorf("Unable to install in namespace: %v", err.Error())
	}

	return namespacedPath, nil
}

func (p *PctInstaller) readConfig(configFile string) (info PuppetContentTemplateInfo, err error) {
	fileBytes, err := p.AFS.ReadFile(configFile)
	if err != nil {
		return info, err
	}

	err = yaml.Unmarshal(fileBytes, &info)
	return info, err
}

func (p *PctInstaller) setupTemplateNamespace(targetDir string, info PuppetContentTemplateInfo, untarPath string, force bool) (string, error) {
	// author/id/version
	templatePath := filepath.Join(targetDir, info.Template.Author, info.Template.Id)

	err := p.AFS.MkdirAll(templatePath, 0750)
	if err != nil {
		return "", err
	}

	namespacePath := filepath.Join(targetDir, info.Template.Author, info.Template.Id, info.Template.Version)

	// finally move to the full path
	err = p.AFS.Rename(untarPath, namespacePath)
	// unable to check for this specific error as windows may instead return `access denied`
	// if err != nil && os.IsExist(err) {
	if err != nil {
		// if a template already exists
		if !force {
			// error unless forced
			return "", fmt.Errorf("Template already installed (%s)", namespacePath)
		} else {
			// remove the exiting template
			err = p.AFS.RemoveAll(namespacePath)
			if err != nil {
				return "", fmt.Errorf("Unable to overwrite existing template: %v", err)
			}
			// perform the move again
			err = p.AFS.Rename(untarPath, namespacePath)
			if err != nil {
				return "", fmt.Errorf("Unable to force install: %v", err)
			}
		}
	}

	return namespacePath, err
}
