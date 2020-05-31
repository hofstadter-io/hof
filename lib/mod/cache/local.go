package cache

import (
	"os"
	"path/filepath"
)

var LocalCacheBaseDir = ".hof/mods"

func init() {
	d, err := os.UserConfigDir()
	if err != nil {
		return
	}

	// save to hof dir for cache across projects
	LocalCacheBaseDir = filepath.Join(d, "hof/mods")
}

func SetBaseDir(basedir string) {
	LocalCacheBaseDir = basedir
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
