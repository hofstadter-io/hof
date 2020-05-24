package resources

import (
	"os"
	"path/filepath"

	"github.com/hofstadter-io/hof/cmd/hof/pflags"
)


var (
	HOF_RESOURCE_DIR = "resources"
	HOF_RESOURCE_BASEDIR = ""
)

func init() {

	if pflags.RootResourcesDirPflag != "" {
		HOF_RESOURCE_BASEDIR = pflags.RootResourcesDirPflag
	} else if bd := os.Getenv("HOF_RESOURCE_BASEDIR"); bd != "" {
		HOF_RESOURCE_BASEDIR = bd
	} else if pflags.RootLocalPflag {
		HOF_RESOURCE_BASEDIR = ".hof"
	} else if pflags.RootGlobalPflag {
		dir, err := os.UserConfigDir()
		if err != nil {
			// return "", err
		}
		// no '.' for a hidden file here
		HOF_RESOURCE_BASEDIR = filepath.Join(dir, "hof")
	} else {
		HOF_RESOURCE_BASEDIR = ".hof"
	}

}
