package mock

import (
	"fmt"
)

type Builder struct {
	ProjectName       string
	ExpectedSourceDir string
	ExpectedTargetDir string
}

func (b *Builder) Build(sourceDir, targetDir string) (gzipArchiveFilePath string, err error) {
	// if input isn't what's expected, raise an error
	if sourceDir != b.ExpectedSourceDir {
		return "", fmt.Errorf("Expected source dir '%s' but got '%s'", b.ExpectedSourceDir, sourceDir)
	}
	if targetDir != b.ExpectedTargetDir {
		return "", fmt.Errorf("Expected source dir '%s' but got '%s'", b.ExpectedSourceDir, sourceDir)
	}
	// If nothing goes wrong, return the path to the packaged project
	return fmt.Sprintf("%s/my-project.tar.gz", targetDir), nil
}
