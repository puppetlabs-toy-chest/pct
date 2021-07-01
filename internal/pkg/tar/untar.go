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
)

func Untar(tarball, target string) (outputDirPath string, err error) {
	reader, err := os.Open(filepath.Clean(tarball))
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
			if err = os.MkdirAll(filepath.Clean(path), info.Mode()); err != nil {
				return "", err
			}
			continue
		}

		file, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())

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
