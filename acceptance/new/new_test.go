package new_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pct/acceptance/testutils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
)

const APP = "pct"

var templatePath string

func TestMain(m *testing.M) {
	log.Logger = zerolog.New(ioutil.Discard).With().Timestamp().Logger()

	templatePath, _ = filepath.Abs("../../internal/pkg/pct/testdata/examples")

	os.Exit(m.Run())
}

func TestPctNew(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	stdout, stderr, exitCode := testutils.RunAppCommand("new", "")
	assert.Contains(t, stdout, "DISPLAYNAME")
	assert.Contains(t, stdout, "AUTHOR")
	assert.Contains(t, stdout, "NAME")
	assert.Contains(t, stdout, "TYPE")
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}

func TestPctNewUnknownTag(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	stdout, stderr, exitCode := testutils.RunAppCommand("new --foo", "")
	assert.Contains(t, stdout, "unknown flag: --foo")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func TestPctNewTemplatePath(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	stdout, stderr, exitCode := testutils.RunAppCommand("new --templatepath "+templatePath, "")
	assert.Contains(t, stdout, "DISPLAYNAME")
	assert.Contains(t, stdout, "NAME")
	assert.Contains(t, stdout, "TYPE")
	assert.Contains(t, stdout, "full-project")
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}

func TestPctNewUnknownTemplate(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	stdout, stderr, exitCode := testutils.RunAppCommand("new foo/bar", "")
	assert.Contains(t, stdout, "Error: Couldn't find an installed template that matches 'foo/bar'")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func TestPctNewAuthorNoId(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	stdout, stderr, exitCode := testutils.RunAppCommand("new puppetlabs", "")
	assert.Contains(t, stdout, "Error: Selected template must be in AUTHOR/ID format")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func TestPctNewKnownTemplate(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	stdout, stderr, exitCode := testutils.RunAppCommand("new puppetlabs/full-project --templatepath "+templatePath, os.TempDir())
	assert.Contains(t, stdout, "Deployed:")
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}

func TestPctNewInfo(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	stdout, stderr, exitCode := testutils.RunAppCommand("new --info puppetlabs/full-project --templatepath "+templatePath, os.TempDir())

	expectedYaml := `puppet_module:
  license: Apache-2.0
  summary: A New Puppet Module
  version: 0.1.0`

	var output map[string]interface{}
	var expected map[string]interface{}

	err := yaml.Unmarshal([]byte(stdout), &output)
	if err != nil {
		assert.Fail(t, "returned data is not YAML")
	}

	err = yaml.Unmarshal([]byte(expectedYaml), &expected)
	if err != nil {
		assert.Fail(t, "expected data is not YAML")
	}

	assert.Equal(t, expected, output)
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}

func TestPctNewInfoJson(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	stdout, stderr, exitCode := testutils.RunAppCommand("new --info puppetlabs/full-project --format json --templatepath "+templatePath, os.TempDir())

	expectedJson := `{
  "puppet_module": {
    "license": "Apache-2.0",
    "version": "0.1.0",
    "summary": "A New Puppet Module"
  }
}`

	var output map[string]interface{}
	var expected map[string]interface{}

	err := json.Unmarshal([]byte(stdout), &output)
	if err != nil {
		assert.Fail(t, "returned data is not JSON")
	}

	err = json.Unmarshal([]byte(expectedJson), &expected)
	if err != nil {
		assert.Fail(t, "expected data is not JSON")
	}

	assert.Equal(t, expected, output)
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}
