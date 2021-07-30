package tar

type TarHelpers interface {
	Tar(source, target string) (tarFilePath string, err error)
}

type TarHelpersImpl struct{}

func (TarHelpersImpl) Tar(source, target string) (tarFilePath string, err error) {
	return Tar(source, target)
}
