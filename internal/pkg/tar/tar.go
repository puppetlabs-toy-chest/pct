package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"strings"

	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type TarI interface {
	Tar(source, target string) (tarFilePath string, err error)
	Untar(tarball, target string) (outputDirPath string, err error)
}

type Tar struct {
	AFS *afero.Afero
}

func (t *Tar) Tar(source, target string) (tarFilePath string, err error) {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarfile, err := t.AFS.Create(target)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := tarfile.Close(); err != nil {
			log.Error().Msgf("Error closing file: %s\n", err)
		}
	}()

	tarball := tar.NewWriter(tarfile)

	defer func() {
		if err := tarball.Close(); err != nil {
			log.Error().Msgf("Error closing writer: %s\n", err)
		}
	}()

	info, err := t.AFS.Stat(source)
	if err != nil {
		return "", nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	err = t.AFS.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := t.AFS.Open(filepath.Clean(path))
			if err != nil {
				return err
			}

			defer func() {
				if err := file.Close(); err != nil {
					log.Error().Msgf("Error closing file: %s\n", err)
				}
			}()

			_, err = io.Copy(tarball, file)
			return err
		})

	if err != nil {
		return "", err
	}

	return target, err
}

func (t *Tar) Untar(tarball, target string) (outputDirPath string, err error) {
	reader, err := t.AFS.Open(filepath.Clean(tarball))
	if err != nil {
		return "", err
	}

	defer func() {
		if err := reader.Close(); err != nil {
			log.Error().Msgf("Error closing reader: %s", err)
		}
	}()

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		if strings.Contains(header.Name, "../") {
			return "", fmt.Errorf("TAR contains entry above top level archive dir: %s\n", header.Name)
		}

		path := filepath.Join(target, filepath.Clean(header.Name))
		info := header.FileInfo()
		if info.IsDir() {
			if err = t.AFS.MkdirAll(filepath.Clean(path), info.Mode()); err != nil {
				return "", err
			}
			continue
		}

		file, err := t.AFS.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())

		defer func() {
			if err := file.Close(); err != nil {
				log.Error().Msgf("Error closing file: %s", err)
			}
		}()

		if err != nil {
			return "", err
		}

		err = utils.ChunkedCopy(file, tarReader)
		if err != nil {
			return "", err
		}

		err = file.Close()
		if err != nil {
			return "", err
		}
	}

	targetDirName := strings.Split(filepath.Base(tarball), ".")[0]
	targetDirPath := filepath.Join(target, targetDirName)

	return targetDirPath, nil
}
