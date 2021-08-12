package mock

type IoUtil struct {
	TestDir string
}

func (i *IoUtil) TempDir() (name string, err error) {
	return i.TestDir, nil
}
