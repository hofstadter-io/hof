package cache

import (
	"path/filepath"

	"github.com/hofstadter-io/hof/lib/mod/util"
)

var (
	LocalCacheBaseDir = filepath.Join(util.UserHomeDir(), ".mvs")
)

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
