package mod

import (
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"

	"github.com/hofstadter-io/hof/lib/mod/mvs"
	"github.com/hofstadter-io/hof/lib/repos/cache"
)

func (cm *CueMod) SolveMVS() error {
	// fmt.Println("solve mvs:", cm.Module)

	rr := &RequirementResolver{ cm }

	targets := []module.Version{}
	for _, dep := range cm.Replace {
		dep := module.Version{ Path: dep.Path, Version: dep.Version }
		targets = append(targets, dep)
	}
	for path, ver := range cm.Require {
		// skip any replaced modules
		if _, ok := cm.Replace[path]; ok {
			continue
		}
		dep := module.Version{ Path: path, Version: ver }
		targets = append(targets, dep)
	}

	final, err := mvs.BuildList(targets, rr)
	if err != nil {
		return err
	}

	cm.BuildList = final

	return nil
}

// type needed by Go's mvs library
type RequirementResolver struct {
	rootMod *CueMod
}

func cmpVersion(v1, v2 string) int {
	if v2 == "" {
		if v1 == "" {
			return 0
		}
		return -1
	}
	if v1 == "" {
		return 1
	}
	return semver.Compare(v1, v2)
}

func (rr *RequirementResolver) Max(v1, v2 string) string {
	if cmpVersion(v1, v2) < 0 {
		return v2
	}
	return v1

	return v1
}

func (rr *RequirementResolver) Required(m module.Version) ([]module.Version, error) {

	var FS billy.Filesystem
	var err error

	// assume local module, can only happen from root module we are working on
	if m.Version == "" {
		// a REPLACE with relative path from the root module
		FS = osfs.New(filepath.Join(rr.rootMod.Basedir, m.Path))
	} else {
		// a REMOTE module
		FS, err = cache.Load(m.Path, m.Version)
	}
	if err != nil {
		return nil, err
	}

	M, err := ReadModule("", FS)
	if err != nil {
		return nil, err
	}

	deps := []module.Version{}
	for path, ver := range M.Require {
		// filter a replaced module or same module?
		if _, ok := rr.rootMod.Replace[path]; ok || path == rr.rootMod.Module {
			continue
		}
		dep := module.Version{ Path: path, Version: ver }
		deps = append(deps, dep)
	}

	return deps, nil
}

/*
	TODO, implement these via the cache
func (rr *RequirementResolver) Upgrade(m module.Version) (module.Version, error) {

	return m, nil
}

func (rr *RequirementResolver) Previous(m module.Version) (module.Version, error) {

	return m, nil
}
*/

