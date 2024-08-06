/*
This package controls all interactions with the config directory and everything therein.
cfgDirPath is guaranteed to be set (init panics on failure) prior to `main()`.
Also handles token path.
*/
package cfgdir

import (
	"os"
	"path"
)

const (
	SubDirName       string = "revolttui"
	defaultTokenName string = "token"
)

// path to current config directory
var cfgDirPath string

// on boot, determine a config directory to use.
// Try the default config directory, but fall back to the local directory on failure.
func init() {
	// attempt to operate out of the default config directory
	var cfgDirErr, pwdDirErr error
	if cfgDirPath, cfgDirErr = trySpawnDir(""); cfgDirErr != nil {
		// failed, try local directory
		if cfgDirPath, pwdDirErr = trySpawnDir("."); pwdDirErr != nil {
			panic(
				"failed to create a suitable config directory.\n" +
					"cfgdir failure: " + cfgDirErr.Error() + "\n" +
					"pwddir failure: " + pwdDirErr.Error(),
			)
		}
	}
	// if we reached this point, we have a usable config directory
	if cfgDirPath == "" { // sanity check
		panic("spawned a config directory, but cfgDirPath is unset")
	}
}

// helper function for initializing the given folder as the config directory.
// If a path is not given, defaults to trying the config directory.
func trySpawnDir(base string) (subDirPath string, err error) {
	if base == "" {
		base, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
	}
	// create the subdirectory
	subdir := path.Join(base, SubDirName)
	if err := os.MkdirAll(subdir, 0777); err != nil {
		return "", err
	}

	return subdir, nil
}

// Returns the path to our config directory. Guaranteed to exist.
func Get() string {
	return cfgDirPath
}
