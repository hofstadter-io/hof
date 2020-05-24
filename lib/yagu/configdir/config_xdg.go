// +build !windows,!darwin

package configdir

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	systemConfig []string
	localConfig  string
	localCache   string
)

func findPaths() {
	// System-wide configuration.
	if os.Getenv("XDG_CONFIG_DIRS") != "" {
		systemConfig = strings.Split(os.Getenv("XDG_CONFIG_DIRS"), ":")
	} else {
		systemConfig = []string{"/etc/xdg"}
	}

	// Local user configuration.
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		localConfig = os.Getenv("XDG_CONFIG_HOME")
	} else {
		localConfig = filepath.Join(os.Getenv("HOME"), ".config")
	}

	// Local user cache.
	if os.Getenv("XDG_CACHE_HOME") != "" {
		localCache = os.Getenv("XDG_CACHE_HOME")
	} else {
		localCache = filepath.Join(os.Getenv("HOME"), ".cache")
	}
}
