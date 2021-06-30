package gzip

import (
	"compress/gzip"
	"os"
	"path/filepath"

	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog/log"
)

func Gunzip(source, target string) error {
	reader, err := os.Open(filepath.Clean(source))
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
	writer, err := os.Create(target)
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
