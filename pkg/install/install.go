package install

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/puppetlabs/pdkgo/pkg/exec_runner"

	"github.com/puppetlabs/pdkgo/pkg/gzip"
	"github.com/puppetlabs/pdkgo/pkg/httpclient"
	"github.com/puppetlabs/pdkgo/pkg/tar"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type ConfigParams struct {
	Id      string `mapstructure:"id"`
	Author  string `mapstructure:"author"`
	Version string `mapstructure:"version"`
}

type Installer struct {
	Tar           tar.TarI
	Gunzip        gzip.GunzipI
	AFS           *afero.Afero
	IOFS          *afero.IOFS
	HTTPClient    httpclient.HTTPClientI
	Exec          exec_runner.ExecI
	InstallConfig InstallConfigI
}

type InstallerI interface {
	Install(templatePkg, targetDir string, force bool) (string, error)
	InstallClone(gitUri, targetDir, tempDir string, force bool) (string, error)
}

type InstallConfigI interface {
	ProcessConfig(sourceDir, targetDir string, force bool) (string, error)
}

func (p *Installer) Install(templatePkg, targetDir string, force bool) (string, error) {
	// Check if the template package path is a url
	if strings.HasPrefix(templatePkg, "http") {
		// Download the tar.gz file and change templatePkg to its download path
		err := p.processDownload(&templatePkg)
		if err != nil {
			return "", err
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

	// Process the configuration file and set up namespacedPath and relocate config and content to it
	namespacedPath, err := p.InstallConfig.ProcessConfig(untarPath, targetDir, force)
	if err != nil {
		return "", fmt.Errorf("Invalid config: %v", err.Error())
	}

	return namespacedPath, nil
}

func (p *Installer) processDownload(templatePkg *string) error {
	u, err := url.ParseRequestURI(*templatePkg)
	if err != nil {
		return fmt.Errorf("Could not parse template url %s: %v", *templatePkg, err)
	}
	// Create a temporary Directory to download the tar.gz to
	tempDownloadDir, err := p.AFS.TempDir("", "")
	defer func() {
		err := p.AFS.Remove(tempDownloadDir)
		log.Debug().Msgf("Failed to remove temp dir: %v", err)
	}()
	if err != nil {
		return fmt.Errorf("Could not create tempdir to download template: %v", err)
	}
	// Download template and assign location to templatePkg
	*templatePkg, err = p.downloadTemplate(u, tempDownloadDir)
	if err != nil {
		return fmt.Errorf("Could not effectively download template: %v", err)
	}
	return nil
}

func (p *Installer) InstallClone(gitUri, targetDir, tempDir string, force bool) (string, error) {
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

	namespacedPath, err := p.InstallConfig.ProcessConfig(folderPath, targetDir, force)
	if err != nil {
		return "", err
	}

	return namespacedPath, nil
}

func (p *Installer) cloneTemplate(gitUri string, tempDir string) (string, error) {
	clonePath := filepath.Join(tempDir, "temp")
	err := p.Exec.Command("git", "clone", gitUri, clonePath)
	if err != nil {
		return "", err
	}

	output, err := p.Exec.Output()
	log.Info().Msgf(string(output))
	if err != nil {
		return "", err
	}
	return clonePath, nil
}

func (p *Installer) downloadTemplate(targetURL *url.URL, downloadDir string) (string, error) {
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
