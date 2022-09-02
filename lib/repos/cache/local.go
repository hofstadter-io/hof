package cache

import (
	"os"
	"path/filepath"
)

// TODO move this config to hofmod-cli
// maybe to lib/repos or lib/config for now

const CustomCacheBaseDirVar = "HOF_CACHE"
var cacheBaseDir  string

func init() {
	e := os.Getenv(CustomCacheBaseDirVar)
	if e != "" {
		cacheBaseDir = e
	} else {
		d, err := os.UserCacheDir()
		if err != nil {
			return
		}

		// save to hof dir for cache across projects
		cacheBaseDir = filepath.Join(d, "hof/mods")
	}
}

func SetBaseDir(basedir string) {
	cacheBaseDir = basedir
}

/*
func Lookup(modFile, mdr, mod, ver string) (zdata []byte, err error) {

	dir := filepath.Join(
		LocalCacheBaseDir,
		"mod",
		mdr,
		mod,
		"@",
		ver,
	)

	fmt.Println("Cache Lookup:", dir)

	// TODO, lookup

	var buf bytes.Buffer

	m := module.Version{ Path: mod, Version: ver }
	err = zip.CreateFromDir(&buf, m, dir, modFile, []string{"**.*"}, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
*/
