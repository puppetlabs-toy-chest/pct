package tar

type UntarHelpers interface {
	Untar(source, target string) (outputDirPath string, err error)
}

type UntarHelpersImpl struct{}

func (UntarHelpersImpl) Untar(source, target string) (outputDirPath string, err error) {
	return Untar(source, target)
}
