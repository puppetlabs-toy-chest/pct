package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

type UtilsHelperI interface {
	IsModuleRoot() (string, error)
	Dir() (string, error)
}

type UtilsHelper struct{}

// Check if we're currently in the module root dir.
// Return the sanitized file path if we are in a module root, otherwise an empty string.
func (u *UtilsHelper) IsModuleRoot() (string, error) {
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

func (u *UtilsHelper) Dir() (string, error) {
	return homedir.Dir()
}
