package external_validator

import (
	"regexp"
)

type ExternalValidator interface {

	// Set the path to the validator executable and any options we want to pass to it.
	// Will validate the path to the executable.
	// Returns an error: nil if the path to exec is found, an error if it is not found.
	SetCommand(command string, options []string) (err error)

	// Returns true if the validator can be run in this context, otherwise false.
	ValidInContext() (valid bool)

	// Create a list of targets (i.e. files) to validate based on the arguments passed in.
	// The args will be processed in order:
	// - targets []string iterated over to identify absolute paths to files and directories
	// - Any files within directories that match the includePattern are appended to the list
	// - The list is iterated over and any files that match excludePattern are dropped out
	// CONSIDERATION: Arg list could be too long for certain contexts / OSs
	GenerateTargets(targets []string, includePattern regexp.Regexp, excludePattern regexp.Regexp) (err error)

	// Returns the list of targets (i.e. files) that will be validated. Empty list if none found.
	GetTargetsList() (targets []string)

	// Runs the validator and returns the return code from the app exec
	RunValidator() (retCode int, err error)
}
