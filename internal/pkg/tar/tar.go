package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func Tar(source, target string) (tarFilePath string, err error) {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarfile, err := os.Create(target)
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

	info, err := os.Stat(source)
	if err != nil {
		return "", nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	err = filepath.Walk(source,
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

			file, err := os.Open(filepath.Clean(path))
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
