package gzip

import (
	"compress/gzip"
	"path/filepath"

	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

type GunzipI interface {
	Gunzip(source, target string) (err error)
}

type Gunzip struct {
	AFS *afero.Afero
}

func (g *Gunzip) Gunzip(source, target string) error {
	reader, err := g.AFS.Open(filepath.Clean(source))
	if err != nil {
		return err
	}

	defer func() {
		if err := reader.Close(); err != nil {
			log.Error().Msgf("Error closing reader: %s\n", err)
		}
	}()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	defer func() {
		if err := archive.Close(); err != nil {
			log.Error().Msgf("Error closing reader: %s\n", err)
		}
	}()

	target = filepath.Join(target, archive.Name)
	writer, err := g.AFS.Create(target)
	if err != nil {
		return err
	}

	defer func() {
		if err := writer.Close(); err != nil {
			log.Error().Msgf("Error closing writer: %s\n", err)
		}
	}()

	return utils.ChunkedCopy(writer, archive)
}
