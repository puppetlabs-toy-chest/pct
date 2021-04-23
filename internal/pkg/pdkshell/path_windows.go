// +build windows

package pdkshell

import (
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// getPDKInstallDirectory returns the directory PDK is located in
func getPDKInstallDirectory(shortName bool) (string, error) {
	pdkInstallDir, err := getRegistryStringKey(`SOFTWARE\Puppet Labs\DevelopmentKit`, "RememberedInstallDir64")
	if err != nil {
		return "", err
	}
	if shortName {
		pdkInstallDir, err = getShortPath(pdkInstallDir)
		if err != nil {
			return "", err
		}
	}

	return pdkInstallDir, nil
}

// getRegistryStringKey returns the string value of a specified registry key
func getRegistryStringKey(path string, key string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue(key)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return s, nil
}

// getShortPath returns the 8.3 shortened version of a long path
// It is possible to have access to a file or directory but not have access to
// some of the parent directories of that file or directory. As a result,
// GetShortPathName may fail when it is unable to query the parent directory of a
// path component to determine the short name for that component.
func getShortPath(longPath string) (string, error) {
	uLongPath, err := windows.UTF16PtrFromString(longPath)
	if err != nil {
		return "", err
	}

	length, err := windows.GetShortPathName(uLongPath, nil, 0)
	if err != nil {
		return "", err
	}

	var short uint16
	_, err = windows.GetShortPathName(uLongPath, &short, length)
	if err != nil {
		return "", err
	}

	s := windows.UTF16PtrToString(&short)
	return s, nil
}
