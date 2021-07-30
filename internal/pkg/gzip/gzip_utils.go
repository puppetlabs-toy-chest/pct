package gzip

type GzipHelpers interface {
	Gzip(source, target string) (gzipFilePath string, err error)
}

type GzipHelpersImpl struct{}

func (GzipHelpersImpl) Gzip(source, target string) (gzipFilePath string, err error) {
	return Gzip(source, target)
}
