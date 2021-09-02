package mock

import (
	"fmt"
	"path/filepath"
)

type PctInstaller struct {
	ExpectedTemplatePkg string
	ExpectedTargetDir   string
}

func (p *PctInstaller) Install(templatePkg string, targetDir string, force bool) (string, error) {
	if templatePkg != p.ExpectedTemplatePkg {
		return "", fmt.Errorf("templatePkg (%v) did not match expected value (%v)", templatePkg, p.ExpectedTemplatePkg)
	}

	if targetDir != p.ExpectedTargetDir {
		return "", fmt.Errorf("targetDir (%v) did not match expected value (%v)", targetDir, p.ExpectedTargetDir)
	}

	return filepath.Clean("/unit/test/path"), nil
}
