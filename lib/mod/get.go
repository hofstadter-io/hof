package mod

import (
	"fmt"
	"strings"

	gomod "golang.org/x/mod/module"
	"golang.org/x/mod/semver"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/repos/cache"
)

func Get(module string, rflags flags.RootPflagpole, gflags flags.Mod__GetFlagpole) (error) {
	upgradeHofMods()

	parts := strings.Split(module, "@")
	if len(parts) != 2 {
		return fmt.Errorf("bad module format %q, should be 'domain.com/path/...@version'", module)
	}

	path, ver := parts[0], parts[1]

	cm, err := loadRootMod()
	if err != nil {
		return err
	}

	// check special condition
	updateMvs := false
	if path == "all" && ver == "latest" {
		err = updateAll(cm, rflags, gflags)
		updateMvs = true
	} else {
		err = updateOne(cm, path, ver, rflags, gflags)
	}
	if err != nil {
		return err
	}

	fns := []func () error {
		func () error { return cm.SolveMVS(updateMvs) },
		cm.CleanDeps,
		cm.CleanSums,
		cm.WriteModule,
	}

	for _, fn := range fns {
		err := fn()
		if err != nil {
			return err
		}
	}

	// TODO, figure out link / vendor style
	return cm.Vendor("", rflags.Verbosity)
}

func updateOne(cm *CueMod, path, ver string, rflags flags.RootPflagpole, gflags flags.Mod__GetFlagpole) (error) {
	err := gomod.CheckPath(path)
	if err != nil {
		return fmt.Errorf("bad module name %q, should have domain format 'domain.com/...'", path)
	}

	if path == cm.Module {
		return fmt.Errorf("cannot get current module")
	}

	// check for indirect and delete
	_, ok := cm.Indirect[path]
	if ok {
		delete(cm.Indirect, path)
	}

	// if latest, update
	if ver == "latest" {
		// make sure we have the latest
		_, err = cache.FetchRepoSource(path, "")
		if err != nil {
			return err
		}
		ver, err = cache.GetLatestTag(path, gflags.Prerelease)
		if err != nil {
			return err
		}
		fmt.Println("found:", ver)
	}

	// check for already required at a version equal or greater (no downgrades with get)
	currVer, ok := cm.Require[path]
	if ok {
		if semver.Compare(currVer, ver) >= 0 {
			return fmt.Errorf("%s@%s is already required", path, currVer)
		}
	}

	// add dep to required
	cm.Require[path] = ver


	return nil
}


func updateAll(cm *CueMod, rflags flags.RootPflagpole, gflags flags.Mod__GetFlagpole) (err error) {
	for path, ver := range cm.Require {
		// make sure we have the latest
		_, err = cache.FetchRepoSource(path, "")
		if err != nil {
			return err
		}
		nver, err := cache.GetLatestTag(path, false)
		if err != nil {
			return err
		}

		// fmt.Println(" ur:", path, ver, nver)
		// only update if newer, incase we have specific prereleases
		if semver.Compare(nver, ver) > 1 {
			cm.Require[path] = nver
		}
	}

	for path, ver := range cm.Indirect {
		// make sure we have the latest
		_, err = cache.FetchRepoSource(path, "")
		if err != nil {
			return err
		}
		// explicitely not doing prerelease here, should we allow updating all to pre-releases?
		nver, err := cache.GetLatestTag(path, false)
		if err != nil {
			return err
		}
		// fmt.Println(" ui:", path, ver)

		if semver.Compare(nver, ver) > 1 {
			cm.Indirect[path] = ver
		}
	}

	return nil
}
