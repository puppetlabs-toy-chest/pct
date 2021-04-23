package utils

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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

func GetListOfFlags(cmd *cobra.Command, argsV []string, flagsToIgnore []string) []string {
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
