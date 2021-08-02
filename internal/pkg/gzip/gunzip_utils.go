package gzip

type GunzipHelpers interface {
	Gunzip(source, target string) (err error)
}

type GunzipHelpersImpl struct{}

func (GunzipHelpersImpl) Gunzip(source, target string) (err error) {
	return Gunzip(source, target)
}
