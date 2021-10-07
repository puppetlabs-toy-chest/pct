package pct

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/puppetlabs/pdkgo/internal/pkg/exec_runner"

	"github.com/puppetlabs/pdkgo/internal/pkg/gzip"
	"github.com/puppetlabs/pdkgo/internal/pkg/httpclient"
	"github.com/puppetlabs/pdkgo/internal/pkg/tar"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type PctInstaller struct {
	Tar        tar.TarI
	Gunzip     gzip.GunzipI
	AFS        *afero.Afero
	IOFS       *afero.IOFS
	HTTPClient httpclient.HTTPClientI
	Exec       exec_runner.ExecI
}

type PctInstallerI interface {
	Install(templatePkg string, targetDir string, force bool) (string, error)
	InstallClone(gitUri string, targetDir string, tempDir string, force bool) (string, error)
}

func (p *PctInstaller) Install(templatePkg string, targetDir string, force bool) (string, error) {

	// If the package path is a URI, download tar to temp folder

	if strings.HasPrefix(templatePkg, "http") {
		u, err := url.ParseRequestURI(templatePkg)
		if err != nil {
			return "", fmt.Errorf("Could not parse template url %s: %v", templatePkg, err)
		}
		// create a temporary Directory to download the tar.gz to
		tempDownloadDir, err := p.AFS.TempDir("", "")
		defer func() {
			err := p.AFS.Remove(tempDownloadDir)
			log.Debug().Msgf("Failed to remove temp dir: %v", err)
		}()
		if err != nil {
			return "", fmt.Errorf("Could not create tempdir to download template: %v", err)
		}
		templatePkg, err = p.downloadTemplate(u, tempDownloadDir)
		if err != nil {
			return "", fmt.Errorf("Could not effectively download template: %v", err)
		}
	}

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

func (p *PctInstaller) InstallClone(gitUri string, targetDir string, tempDir string, force bool) (string, error) {
	// Validate git URI
	_, err := url.ParseRequestURI(gitUri)
	if err != nil {
		return "", fmt.Errorf("Could not parse template uri %s: %v", gitUri, err)
	}

	// Clone git repository to temp folder
	folderPath, err := p.cloneTemplate(gitUri, tempDir)
	if err != nil {
		return "", fmt.Errorf("Could not clone git repository: %v", err)
	}

	// Remove .git folder from cloned repository
	err = p.AFS.RemoveAll(filepath.Join(folderPath, ".git"))
	if err != nil {
		return "", fmt.Errorf("Failed to remove '.git' directory")
	}

	// Read config to determine template properties
	info, err := p.readConfig(filepath.Join(folderPath, "pct-config.yml"))
	if err != nil {
		return "", fmt.Errorf("Invalid config: %v", err.Error())
	}

	// Create namespaced directory and move contents of temp folder to it
	namespacedPath, err := p.setupTemplateNamespace(targetDir, info, folderPath, force)
	if err != nil {
		return "", fmt.Errorf("Unable to install in namespace: %v", err.Error())
	}

	return namespacedPath, nil
}

func (p *PctInstaller) cloneTemplate(gitUri string, tempDir string) (string, error) {
	// TODO: Sanitize command args
	//clonePath := filepath.Clean(filepath.Join(tempDir, "temp"))
	clonePath := filepath.Join(tempDir, "temp")
	command := p.Exec.Command("git", "clone", gitUri, clonePath)
	output, err := command.Output()
	log.Info().Msgf(string(output))
	if err != nil {
		return "", err
	}
	return clonePath, nil
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

func (p *PctInstaller) downloadTemplate(targetURL *url.URL, downloadDir string) (string, error) {
	// Get the file contents from URL
	response, err := p.HTTPClient.Get(targetURL.String())
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		message := fmt.Sprintf("Received response code %d when trying to download from %s", response.StatusCode, targetURL.String())
		return "", errors.New(message)
	}

	// Create the empty file
	fileName := filepath.Base(targetURL.Path)
	downloadPath := filepath.Join(downloadDir, fileName)
	file, err := p.AFS.Create(downloadPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write file contents
	err = p.AFS.WriteReader(downloadPath, response.Body)
	if err != nil {
		return "", err
	}

	return downloadPath, nil
}
