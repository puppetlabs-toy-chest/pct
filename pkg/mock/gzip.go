package mock

import (
	"fmt"
)

type Gzip struct {
	ReturnedPath string
	ErrResponse  bool
}

func (g *Gzip) Gzip(source, target string) (gzipFilePath string, err error) {
	if g.ErrResponse {
		return "", fmt.Errorf("gzip error")
	}

	return g.ReturnedPath, nil
}
