package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/puppetlabs/pdkgo/internal/pkg/pdkshell"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type UtilsHelper interface {
	IsModuleRoot() (string, error)
}

type UtilsHelperImpl struct{}

// contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// finds a string present in a slice
func Find(source []string, match string) []string {
	var matches []string
	if Contains(source, match) {
		matches = append(matches, match)
	}
	return matches
}

// GetListOfFlags returns a filtered list of arguments provided by the user,
// removing the flags that are not used by PDK Ruby
func GetListOfFlags(cmd *cobra.Command, argsV []string) []string {
	flagsToIgnore := FlagsToIgnore()
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !Contains(flagsToIgnore, f.Name) {
			if f.Changed {
				switch f.Value.Type() {
				case "bool":
					argsV = append(argsV, fmt.Sprintf("--%v", f.Name))
				case "string":
					argsV = append(argsV, fmt.Sprintf("--%v %v", f.Name, f.Value))
				}
			}
		}
	})
	return argsV
}

// FlagsToIgnore list of pdkgo flags not for use in pdk ruby
func FlagsToIgnore() []string {
	flagsToIgnore := []string{"log-level"}
	return flagsToIgnore
}

// ExecutePDKCommand is a helper for executing the pdk commandline
func ExecutePDKCommand(cmd *cobra.Command, args []string) error {
	argsV := buildPDKCommandName(cmd)

	argsV = append(argsV, args...)

	argsV = GetListOfFlags(cmd, argsV)

	log.Trace().Msgf("args: %v", argsV)

	_, err := pdkshell.Execute(argsV)

	return err
}

func buildPDKCommandName(cmd *cobra.Command) []string {
	var argsV []string
	if cmd.HasParent() && cmd.Parent().Name() != "pct" {
		argsV = append(argsV, cmd.Parent().Name(), cmd.Name())
	} else {
		argsV = append(argsV, cmd.Name())
	}
	return argsV
}

// Implement a safe way to copy from src -> dst to prevent
// cases of DoS via decompression bombs (CWE-409)
const copyWriteSize = 1024
const maxCopySize = copyWriteSize * 1024 * 10 // 10MB copy limit

func ChunkedCopy(dst io.Writer, src io.Reader) error {
	currentBytesWritten := int64(0)
	for {
		n, err := io.CopyN(dst, src, copyWriteSize)
		currentBytesWritten += n
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if currentBytesWritten > maxCopySize {
			return fmt.Errorf("Exceeded max copy size of %v", maxCopySize)
		}
	}
}

// Check if we're currently in the module root dir.
// Return the sanitized file path if we are in a module root, otherwise an empty string.
func (UtilsHelperImpl) IsModuleRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	_, err = os.Stat(filepath.Join(wd, "metadata.json"))
	if err != nil {
		return "", err
	}

	return filepath.Clean(wd), nil
}
