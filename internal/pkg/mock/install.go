package mock

import (
	"fmt"
	"path/filepath"
)

type PctInstaller struct {
	ExpectedTemplatePkg string
	ExpectedTargetDir   string
	ExpectedGitUri      string
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

func (p *PctInstaller) InstallClone(gitUri, targetDir, tempDir string, force bool) (string, error) {
	if gitUri != p.ExpectedGitUri {
		return "", fmt.Errorf("gitUri (%v) did not match expected value (%v)", gitUri, p.ExpectedGitUri)
	}

	if tempDir == "" {
		return "", fmt.Errorf("tempDir was an empty string")
	}

	if targetDir != p.ExpectedTargetDir {
		return "", fmt.Errorf("targetDir (%v) did not match expected value (%v)", targetDir, p.ExpectedTargetDir)
	}

	return filepath.Clean("/unit/test/path"), nil
}
