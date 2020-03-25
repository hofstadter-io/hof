package extern

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hofstadter-io/hof/lib/util"
)

func ImportAddBundle(bundle string) (string, error) {
	if bundle == "" || bundle[0] == '#' || bundle[0] == '@' {
		bundle = "https://github.com/hofstadter-io/studios-modules" + bundle
	}
	url, version, subpath := SplitParts(bundle)

	err := cloneAndRenderImport(url, version, subpath)
	if err != nil {
		return "", err
	}

	// TODO update some deps file

	return "Done", nil
}

func cloneAndRenderImport(srcUrl, srcVer, srcPath string) error {
	_, appname := util.GetAcctAndName()
	data := map[string]interface{}{
		"AppName": appname,
	}

	dir, err := util.CloneRepo(srcUrl, srcVer)
	if err != nil {
		return err
	}

	err = util.RenderDir(filepath.Join(dir, srcPath, "design"), "design-vendor", data)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(dir, srcPath, "design-vendor")); !os.IsNotExist(err) {
		// path exists
		err = util.RenderDir(filepath.Join(dir, srcPath, "design-vendor"), "design-vendor", data)
		if err != nil {
			return err
		}
	}
	return nil
}

func SplitParts(full string) (url, version, subpath string) {
	posVersion := strings.LastIndex(full, "@")
	posSubpath := strings.LastIndex(full, "#")

	if posVersion == -1 && posSubpath == -1 {
		url = full
		return
	}

	if posVersion == -1 {
		parts := strings.Split(full, "#")
		url, subpath = parts[0], parts[1]
		return
	}

	if posSubpath == -1 {
		parts := strings.Split(full, "@")
		url, version = parts[0], parts[1]
		return
	}

	if posVersion < posSubpath {
		parts := strings.Split(full, "#")
		subpath = parts[1]
		parts = strings.Split(parts[0], "@")
		url, version = parts[0], parts[1]
	} else {
		parts := strings.Split(full, "@")
		version = parts[1]
		parts = strings.Split(parts[0], "#")
		url, subpath = parts[0], parts[1]
	}

	return url, version, subpath
}
