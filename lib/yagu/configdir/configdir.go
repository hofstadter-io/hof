package configdir

import (
	"os"
	"path/filepath"
)

// VERSION is the semantic version number of the configdir library.
const VERSION = "0.1.0"

func init() {
	findPaths()
}

// Refresh will rediscover the config paths, checking current environment
// variables again.
//
// This function is automatically called when the program initializes. If you
// change the environment variables at run-time, though, you may call the
// Refresh() function to reevaluate the config paths.
func Refresh() {
	findPaths()
}

// SystemConfig returns the system-wide configuration paths, with optional path
// components added to the end for vendor/application-specific settings.
func SystemConfig(folder ...string) []string {
	if len(folder) == 0 {
		return systemConfig
	}

	var paths []string
	for _, root := range systemConfig {
		p := append([]string{root}, filepath.Join(folder...))
		paths = append(paths, filepath.Join(p...))
	}

	return paths
}

// LocalConfig returns the local user configuration path, with optional
// path components added to the end for vendor/application-specific settings.
func LocalConfig(folder ...string) string {
	if len(folder) == 0 {
		return localConfig
	}

	return filepath.Join(localConfig, filepath.Join(folder...))
}

// LocalCache returns the local user cache folder, with optional path
// components added to the end for vendor/application-specific settings.
func LocalCache(folder ...string) string {
	if len(folder) == 0 {
		return localCache
	}

	return filepath.Join(localCache, filepath.Join(folder...))
}

// DefaultFileMode controls the default permissions on any paths created by
// using MakePath.
var DefaultFileMode = os.FileMode(0755)

// MakePath ensures that the full path you wanted, including vendor or
// application-specific components, exists. You can give this the output of
// any of the config path functions (SystemConfig, LocalConfig or LocalCache).
//
// In the event that the path function gives multiple answers, e.g. for
// SystemConfig, MakePath() will only attempt to create the sub-folders on
// the *first* path found. If this isn't what you want, you may want to just
// use the os.MkdirAll() functionality directly.
func MakePath(paths ...string) error {
	if len(paths) >= 1 {
		err := os.MkdirAll(paths[0], DefaultFileMode)
		if err != nil {
			return err
		}
	}

	return nil
}
