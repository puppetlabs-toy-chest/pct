package gzip

import (
	"compress/gzip"
	"path/filepath"

	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type GunzipI interface {
	Gunzip(source, target string) (string, error)
}

type Gunzip struct {
	AFS *afero.Afero
}

// returns the filepath to the extracted file. If there is an error the string will be empty
func (g *Gunzip) Gunzip(source, target string) (string, error) {
	reader, err := g.AFS.Open(filepath.Clean(source))
	if err != nil {
		return "",err
	}

	defer func() {
		if err := reader.Close(); err != nil {
			log.Error().Msgf("Error closing reader: %s\n", err)
		}
	}()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return "",err
	}

	defer func() {
		if err := archive.Close(); err != nil {
			log.Error().Msgf("Error closing reader: %s\n", err)
		}
	}()

	tar := filepath.Join(target, archive.Name)
	writer, err := g.AFS.Create(tar)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := writer.Close(); err != nil {
			log.Error().Msgf("Error closing writer: %s\n", err)
		}
	}()

	return tar, utils.ChunkedCopy(writer, archive)
}
