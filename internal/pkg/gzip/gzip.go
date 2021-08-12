package gzip

import (
	"compress/gzip"
	"fmt"
	"path/filepath"

	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type GzipI interface {
	Gzip(source, target string) (gzipFilePath string, err error)
}

type Gzip struct {
	AFS *afero.Afero
}

func (g *Gzip) Gzip(source, target string) (gzipFilePath string, err error) {
	reader, err := g.AFS.Open(filepath.Clean(source))
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(target, 0750)
	if err != nil {
		return "", err
	}

	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := g.AFS.Create(target)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := writer.Close(); err != nil {
			log.Error().Msgf("Error closing writer: %s\n", err)
		}
	}()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename

	defer func() {
		if err := archiver.Close(); err != nil {
			log.Error().Msgf("Error closing writer: %s\n", err)
		}
	}()

	err = utils.ChunkedCopy(archiver, reader)
	if err != nil {
		return "", err
	}

	return target, err
}
