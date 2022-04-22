package build_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pct/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

const APP = "pct"

func Test_PctBuild_Outputs_TarGz(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	templateName := "good-project"

	sourceDir, _ := filepath.Abs("../../acceptance/build/testdata")
	templateDir := filepath.Join(sourceDir, templateName)
	wd := testutils.GetTmpDir(t)

	cmd := fmt.Sprintf("build --sourcedir %v --targetdir %v", templateDir, wd)
	stdout, stderr, exitCode := testutils.RunAppCommand(cmd, "")

	expectedOutputFilePath := filepath.Join(wd, fmt.Sprintf("%v.tar.gz", templateName))

	assert.Contains(t, stdout, fmt.Sprintf("Packaged template output to %v", expectedOutputFilePath))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
	assert.FileExists(t, expectedOutputFilePath)
}

func Test_PctBuild_With_NoTargetDir_Outputs_TarGz(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	templateName := "good-project"

	sourceDir, _ := filepath.Abs("../../acceptance/build/testdata")
	templateDir := filepath.Join(sourceDir, templateName)
	wd := testutils.GetTmpDir(t)

	cmd := fmt.Sprintf("build --sourcedir %v", templateDir)
	stdout, stderr, exitCode := testutils.RunAppCommand(cmd, wd)

	expectedOutputFilePath := filepath.Join(wd, "pkg", fmt.Sprintf("%v.tar.gz", templateName))

	assert.Contains(t, stdout, fmt.Sprintf("Packaged template output to %v", expectedOutputFilePath))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
	assert.FileExists(t, expectedOutputFilePath)
}

func Test_PctBuild_With_EmptySourceDir_Errors(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	templateName := "no-project-here"

	sourceDir, _ := filepath.Abs("../../acceptance/build/testdata")
	templateDir := filepath.Join(sourceDir, templateName)

	cmd := fmt.Sprintf("build --sourcedir %v", templateDir)
	stdout, stderr, exitCode := testutils.RunAppCommand(cmd, "")

	assert.Contains(t, stdout, fmt.Sprintf("No project directory at %v", templateDir))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PctBuild_With_NoPctConfig_Errors(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	templateName := "no-pct-config-project"

	sourceDir, _ := filepath.Abs("../../acceptance/build/testdata")
	templateDir := filepath.Join(sourceDir, templateName)

	cmd := fmt.Sprintf("build --sourcedir %v", templateDir)
	stdout, stderr, exitCode := testutils.RunAppCommand(cmd, "")

	assert.Contains(t, stdout, fmt.Sprintf("No 'pct-config.yml' found in %v", templateDir))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PctBuild_With_NoContentDir_Errors(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	templateName := "no-content-dir-project"

	sourceDir, _ := filepath.Abs("../../acceptance/build/testdata")
	templateDir := filepath.Join(sourceDir, templateName)

	cmd := fmt.Sprintf("build --sourcedir %v", templateDir)
	stdout, stderr, exitCode := testutils.RunAppCommand(cmd, "")

	assert.Contains(t, stdout, fmt.Sprintf("No 'content' dir found in %v", templateDir))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}
