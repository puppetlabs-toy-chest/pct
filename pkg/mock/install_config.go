package mock

import "fmt"

type InstallConfig struct {
	ExpectedSourceDir      string
	ExpectedTargetDir      string
	ExpectedForce          bool
	NamespacedPathResponse string

	ErrResponse error
}

func (ic *InstallConfig) ProcessConfig(sourceDir, targetDir string, force bool) (string, error) {
	if ic.ErrResponse != nil {
		return "", ic.ErrResponse
	}

	if sourceDir != ic.ExpectedSourceDir {
		return "", fmt.Errorf("sourceDir (%v) did not match expected value (%v)", sourceDir, ic.ExpectedSourceDir)
	}

	if targetDir != ic.ExpectedTargetDir {
		return "", fmt.Errorf("targetDir (%v) did not match expected value (%v)", targetDir, ic.ExpectedTargetDir)
	}

	if force != ic.ExpectedForce {
		return "", fmt.Errorf("force (%v) did not match expected value (%v)", force, ic.ExpectedForce)
	}
	return ic.NamespacedPathResponse, nil
}
