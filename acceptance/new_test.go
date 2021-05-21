package new_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pdkgo/acceptance/testutils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/stretchr/testify/assert"
)

var templatePath string;

func TestMain(m *testing.M) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	templatePath, _ = filepath.Abs("../internal/pkg/pct/testdata/examples")

	os.Exit(m.Run())
}

func TestPctNew(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	stdout, stderr, exitCode := testutils.RunPctCommand("new", "")
	assert.Contains(t, stdout, "DISPLAYNAME" )
	assert.Contains(t, stdout, "NAME" )
	assert.Contains(t, stdout, "TYPE" )
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}

func TestPctNewUnknownTag(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	stdout, stderr, exitCode := testutils.RunPctCommand("new --foo", "")
	assert.Contains(t, stdout, "unknown flag: --foo", )
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func TestPctNewTemplatePath(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	stdout, stderr, exitCode := testutils.RunPctCommand("new --templatepath " + templatePath, "")
	assert.Contains(t, stdout, "DISPLAYNAME" )
	assert.Contains(t, stdout, "NAME" )
	assert.Contains(t, stdout, "TYPE" )
	assert.Contains(t, stdout, "full-project" )
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}

func TestPctNewUnknownTemplate(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	stdout, stderr, exitCode := testutils.RunPctCommand("new foo", "")
	assert.Contains(t, stdout, "Error: Couldn't find an installed template that matches 'foo'" )
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func TestPctNewKnownTemplate(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	stdout, stderr, exitCode := testutils.RunPctCommand("new full-project --templatepath " + templatePath, os.TempDir())
	assert.Contains(t, stdout, "Deployed:" )
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}
