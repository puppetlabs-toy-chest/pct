package ruby

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/puppetlabs/pdkgo/internal/pkg/pdkshell"
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	debug         bool
	autoCorrect   bool
	list          bool
	parallel      bool
	puppetDev     bool
	peVersion     string
	puppetVersion string
)

type myUtils interface {
	ValidModuleRoot() (string, error)
}

type RubyRubocopValidator struct{}
type utilHelper struct{
	foo myUtils
}

var commandAndArgs []string

// Assign these functions to variables to mock during tests
var fValidModuleRoot = utils.ValidModuleRoot
var fOsStat = os.Stat

func (r *RubyRubocopValidator) SetCommand(command string, options []string) (err error) {
	// This is hard coded for now. The existing implementation determines where it is being called from:
	// - Module root: $MODULE_ROOT/bin/rubocop

	currentWd, err := fValidModuleRoot()
	if (err != nil) {
		log.Error().Msgf("Error determining whether current working dir is module root: %v", err)
		return err
	}

	commandPath := filepath.Join(currentWd, "bin", command)
	_, err = fOsStat(commandPath)
	if err != nil {
		log.Error().Msgf("Could not stat path (%s): %v", commandPath, err)
		return err
	}

	commandAndArgs = append([]string{commandPath}, options...)
	log.Trace().Msgf("Command and args: %v", commandAndArgs)

	return nil
}

func (r *RubyRubocopValidator) RunValidator() (retCode int, err error) {
	return pdkshell.Execute(commandAndArgs)
}

func (r *RubyRubocopValidator) GenerateTargets(targets []string, includePattern regexp.Regexp, excludePattern regexp.Regexp) (err error) {
	// Determine what are dirs from "targets"
	// Glob all subdirs to enumerate files within that match includePattern
	// Drop out anything that matches excludePattern

	return nil
}

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use: 		"ruby [flags]",
		Short:	"Run Ruby code validators",
		Long: 	`Run Ruby code validators`,
		RunE:		executeValidateRuby,
	}

	return tmp
}

func executeValidateRuby(cmd *cobra.Command, targets []string) error {
	var validator RubyRubocopValidator
	var args []string = []string{"--format", "json"}
	if autoCorrect {
		args = append([]string{"--auto-correct"}, args...)
	}

	validator.SetCommand("rubocop", args)
	_, err := validator.RunValidator()
	return err
}
