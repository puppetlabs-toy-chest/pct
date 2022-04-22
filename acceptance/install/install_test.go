package install_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pct/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

const APP = "pct"

var defaultTemplatePath string

func Test_PctInstall_InstallsTo_DefaultTemplatePath(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Setup
	templatePkgPath, _ := filepath.Abs(fmt.Sprintf("../../acceptance/install/testdata/%v.tar.gz", "good-project"))
	installCmd := fmt.Sprintf("install %v", templatePkgPath)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// Assert
	assert.Contains(t, stdout, fmt.Sprintf("Template installed to %v", filepath.Join(getDefaultTemplatePath(), "gooder", "good-project", "0.1.0")))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
	assert.FileExists(t, filepath.Join(getDefaultTemplatePath(), "gooder", "good-project", "0.1.0", "pct-config.yml"))
	assert.FileExists(t, filepath.Join(getDefaultTemplatePath(), "gooder", "good-project", "0.1.0", "content", "empty.txt"))
	assert.FileExists(t, filepath.Join(getDefaultTemplatePath(), "gooder", "good-project", "0.1.0", "content", "goodfile.txt.tmpl"))

	stdout, stderr, exitCode = testutils.RunAppCommand("new --list", "")
	assert.Regexp(t, "Good\\sProject\\s+\\|\\sgooder\\s+\\|\\sgood-project\\s+\\|\\sproject", stdout)
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)

	// Tear Down
	removeInstalledTemplate(filepath.Join(getDefaultTemplatePath(), "gooder", "good-project", "0.1.0"))
}

type templateData struct {
	name          string
	author        string
	listExpRegex  string
	expectedFiles []string
	gitUri        string
}

func Test_PctInstall_InstallsTo_DefinedTemplatePath(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Setup
	templatePath := testutils.GetTmpDir(t)

	templatePkgs := []templateData{
		{
			name:         "additional-project",
			author:       "adder",
			listExpRegex: "Additional\\sProject\\s+\\|\\sadder\\s+\\|\\sadditional-project\\s+\\|\\sproject",
			expectedFiles: []string{
				"pct-config.yml",
				"content/empty.txt",
				"content/goodfile.txt.tmpl",
			},
		},
		{
			name:         "good-project",
			author:       "gooder",
			listExpRegex: "Good\\sProject\\s+\\|\\sgooder\\s+\\|\\sgood-project\\s+\\|\\sproject",
			expectedFiles: []string{
				"pct-config.yml",
				"content/empty.txt",
				"content/goodfile.txt.tmpl",
			},
		},
	}

	for _, template := range templatePkgs {
		// Setup
		templatePkgPath, _ := filepath.Abs(fmt.Sprintf("../../acceptance/install/testdata/%v.tar.gz", template.name))
		installCmd := fmt.Sprintf("install %v --templatepath %v", templatePkgPath, templatePath)

		// Exec
		stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

		// Assert
		assert.Contains(t, stdout, fmt.Sprintf("Template installed to %v", filepath.Join(templatePath, template.author, template.name, "0.1.0")))
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	for _, template := range templatePkgs {
		// Assert
		for _, file := range template.expectedFiles {
			assert.FileExists(t, filepath.Join(templatePath, template.author, template.name, "0.1.0", file))
		}

		listCmd := fmt.Sprintf("new --list --templatepath %v", templatePath)
		stdout, stderr, exitCode := testutils.RunAppCommand(listCmd, "")

		assert.Regexp(t, template.listExpRegex, stdout)
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	// Tear Down
	for _, template := range templatePkgs {
		removeInstalledTemplate(filepath.Join(templatePath, template.author, template.name, "0.1.0"))
	}
}

func Test_PctInstall_InstallsFrom_RemoteTemplatePath(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Setup
	templatePath := testutils.GetTmpDir(t)

	templatePkgs := []templateData{
		{
			name:         "additional-project",
			author:       "adder",
			listExpRegex: "Additional\\sProject\\s+\\|\\sadder\\s+\\|\\sadditional-project\\s+\\|\\sproject",
			expectedFiles: []string{
				"pct-config.yml",
				"content/empty.txt",
				"content/goodfile.txt.tmpl",
			},
		},
		{
			name:         "good-project",
			author:       "gooder",
			listExpRegex: "Good\\sProject\\s+\\|\\sgooder\\s+\\|\\sgood-project\\s+\\|\\sproject",
			expectedFiles: []string{
				"pct-config.yml",
				"content/empty.txt",
				"content/goodfile.txt.tmpl",
			},
		},
	}

	for _, template := range templatePkgs {
		// Setup
		templatePkgPath := fmt.Sprintf("https://github.com/puppetlabs/pct/raw/main/acceptance/install/testdata/%s.tar.gz", template.name)
		installCmd := fmt.Sprintf("install %v --templatepath %v", templatePkgPath, templatePath)

		// Exec
		stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

		// Assert
		assert.Contains(t, stdout, fmt.Sprintf("Template installed to %v", filepath.Join(templatePath, template.author, template.name, "0.1.0")))
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	for _, template := range templatePkgs {
		// Assert
		for _, file := range template.expectedFiles {
			assert.FileExists(t, filepath.Join(templatePath, template.author, template.name, "0.1.0", file))
		}

		listCmd := fmt.Sprintf("new --list --templatepath %v", templatePath)
		stdout, stderr, exitCode := testutils.RunAppCommand(listCmd, "")

		assert.Regexp(t, template.listExpRegex, stdout)
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	// Tear Down
	for _, template := range templatePkgs {
		removeInstalledTemplate(filepath.Join(templatePath, template.author, template.name, "0.1.0"))
	}
}

func Test_PctInstall_Errors_When_NoTemplatePkgDefined(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("install", "")

	// Assert
	assert.Contains(t, stdout, "Path to template package (tar.gz) should be first argument")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PctInstall_Errors_When_TemplatePkgNotExist(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Setup
	templatePkgPath, _ := filepath.Abs("/path/to/nowhere/good-project.tar.gz")
	installCmd := fmt.Sprintf("install %v", templatePkgPath)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// Assert
	assert.Contains(t, stdout, fmt.Sprintf("No package at %v", templatePkgPath))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PctInstall_Errors_When_InvalidGzProvided(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Setup
	templatePkgPath, _ := filepath.Abs("../../acceptance/install/testdata/invalid-gz-project.tar.gz")
	installCmd := fmt.Sprintf("install %v", templatePkgPath)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// Assert
	assert.Contains(t, stdout, fmt.Sprintf("Could not extract TAR from GZIP (%v)", templatePkgPath))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PctInstall_Errors_When_InvalidTarProvided(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Setup
	templatePkgPath, _ := filepath.Abs("../../acceptance/install/testdata/invalid-tar-project.tar.gz")
	installCmd := fmt.Sprintf("install %v", templatePkgPath)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// Assert
	assert.Contains(t, stdout, fmt.Sprintf("Could not UNTAR package (%v)", templatePkgPath))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PctInstall_FailsWhenTemplateAlreadyExists(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Setup
	templatePath := testutils.GetTmpDir(t)

	// Install template
	templatePkgPath, _ := filepath.Abs("../../acceptance/install/testdata/additional-project.tar.gz")
	installCmd := fmt.Sprintf("install %v --templatepath %v", templatePkgPath, templatePath)
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// verify the template installed
	assert.Contains(t, stdout, fmt.Sprintf("Template installed to %v", filepath.Join(templatePath, "adder", "additional-project", "0.1.0")))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)

	// Attempt to reinstall the template
	templatePkgPath, _ = filepath.Abs("../../acceptance/install/testdata/additional-project.tar.gz")
	installCmd = fmt.Sprintf("install %v --templatepath %v", templatePkgPath, templatePath)
	stdout, stderr, exitCode = testutils.RunAppCommand(installCmd, "")

	// verify that the template failed to install
	assert.Contains(t, stdout, "Unable to install in namespace: Package already installed")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)

	// Tear Down
	removeInstalledTemplate(filepath.Join(templatePath, "adder", "additional-project", "0.1.0"))
}

func Test_PctInstall_ForceSuccessWhenTemplateAlreadyExists(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Setup
	templatePath := testutils.GetTmpDir(t)

	// Install template
	templatePkgPath, _ := filepath.Abs("../../acceptance/install/testdata/additional-project.tar.gz")
	installCmd := fmt.Sprintf("install %v --templatepath %v", templatePkgPath, templatePath)
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// verify the template installed
	assert.Contains(t, stdout, fmt.Sprintf("Template installed to %v", filepath.Join(templatePath, "adder", "additional-project", "0.1.0")))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)

	// Attempt to reinstall the template
	templatePkgPath, _ = filepath.Abs("../../acceptance/install/testdata/additional-project.tar.gz")
	installCmd = fmt.Sprintf("install %v --force --templatepath %v", templatePkgPath, templatePath)
	stdout, stderr, exitCode = testutils.RunAppCommand(installCmd, "")

	// verify that the template reinstall exited successfully
	assert.Contains(t, stdout, fmt.Sprintf("Template installed to %v", filepath.Join(templatePath, "adder", "additional-project", "0.1.0")))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)

	// Tear Down
	removeInstalledTemplate(filepath.Join(templatePath, "adder", "additional-project", "0.1.0"))
}

// Util Functions

func removeInstalledTemplate(templatePath string) {
	_, err := os.Stat(templatePath)
	if err != nil {
		panic(fmt.Sprintf("removeInstalledTemplate(): Could not determine if Package path (%v) exists: %v", templatePath, err))
	}

	os.RemoveAll(templatePath)
	if err != nil {
		panic(fmt.Sprintf("remoteTemplate(): Could not remove %v: %v", templatePath, err))
	}
}

func getDefaultTemplatePath() string {
	if defaultTemplatePath != "" {
		return defaultTemplatePath
	}

	entries, err := filepath.Glob("../../dist/pct_*/templates")
	if err != nil {
		panic("getDefaultTemplatePath(): Could not determine default template path")
	}
	if len(entries) != 1 {
		panic(fmt.Sprintf("getDefaultTemplatePath(): Could not determine default template path; matched entries: %v", len(entries)))
	}

	defaultTemplatePath, err := filepath.Abs(entries[0])
	if err != nil {
		panic(fmt.Sprintf("getDefaultTemplatePath(): Could not create absolute path to templatepath: %v", err))
	}

	return defaultTemplatePath
}

func Test_PctInstall_WithGitUri_InstallTemplate(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Setup
	templatePath := testutils.GetTmpDir(t)

	templatePkgs := []templateData{
		{
			name:         "test-template-1",
			author:       "test-user",
			listExpRegex: "Test\\sTemplate\\s1\\s+\\|\\stest-user\\s+\\|\\stest-template-1\\s+\\|\\sproject",
			expectedFiles: []string{
				"pct-config.yml",
			},
			gitUri: "https://github.com/puppetlabs/pct-test-template-01.git",
		},
		{
			name:         "test-template-2",
			author:       "test-user",
			listExpRegex: "Test\\sTemplate\\s2\\s+\\|\\stest-user\\s+\\|\\stest-template-2\\s+\\|\\sproject",
			expectedFiles: []string{
				"pct-config.yml",
			},
			gitUri: "https://github.com/puppetlabs/pct-test-template-02.git",
		},
	}

	for _, template := range templatePkgs {
		installCmd := fmt.Sprintf("install --git-uri %v --templatepath %v", template.gitUri, templatePath)

		// Exec
		stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

		// Assert
		assert.Contains(t, stdout, fmt.Sprintf("Template installed to %v", filepath.Join(templatePath, template.author, template.name, "0.1.0")))
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	for _, template := range templatePkgs {
		// Assert
		for _, file := range template.expectedFiles {
			assert.FileExists(t, filepath.Join(templatePath, template.author, template.name, "0.1.0", file))
		}

		listCmd := fmt.Sprintf("new --list --templatepath %v", templatePath)
		stdout, stderr, exitCode := testutils.RunAppCommand(listCmd, "")

		assert.Regexp(t, template.listExpRegex, stdout)
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	// Tear Down
	for _, template := range templatePkgs {
		removeInstalledTemplate(filepath.Join(templatePath, template.author, template.name, "0.1.0"))
	}
}

func Test_PctInstall_WithGitUri_FailsWithNonExistentUri(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("install --git-uri https://example.com/fake-git-uri", "")

	// Assert
	assert.Contains(t, stdout, "Could not clone git repository:")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PctInstall_WithGitUri_FailsWithInvalidUri(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("install --git-uri example.com/invalid-git-uri", "")

	// Assert
	assert.Contains(t, stdout, "Could not parse package uri")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PctInstall_WithGitUri_RemovesHiddenGitDir(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	// Setup
	templatePath := testutils.GetTmpDir(t)

	// Install template
	installCmd := fmt.Sprint("install --git-uri https://github.com/puppetlabs/pct-test-template-01.git --templatepath ", templatePath)
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// Verify the template installed
	assert.Contains(t, stdout, fmt.Sprintf("Template installed to %v", filepath.Join(templatePath, "test-user", "test-template-1", "0.1.0")))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)

	// Check .git directory has been deleted
	assert.NoDirExists(t, filepath.Join(templatePath, "test-user", "test-template-1", "0.1.0", ".git"))

	// Tear Down
	removeInstalledTemplate(filepath.Join(templatePath, "test-user", "test-template-1", "0.1.0"))
}
