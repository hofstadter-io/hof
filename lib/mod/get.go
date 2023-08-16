package mod

import (
	"fmt"
	"strings"

	// gomod "golang.org/x/mod/module"
	"golang.org/x/mod/semver"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/repos/cache"
)

func Get(module string, rflags flags.RootPflagpole, gflags flags.Mod__GetFlagpole) (error) {
	upgradeHofMods()

	path, ver := module, ""

	// figure out parts
	parts := []string{module}
	if strings.Contains(module, "@") {
		parts = strings.Split(module, "@")
	} else if strings.Contains(module, ":") {
		parts = strings.Split(module, ":")
	}
	if len(parts) == 2 {
		path, ver = parts[0], parts[1]
	}
	if ver == "" {
		ver = "latest"
	}

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
		cm.UpgradePseudoVersions,
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

func updateOne(cm *CueMod, path, ver string, rflags flags.RootPflagpole, gflags flags.Mod__GetFlagpole) (err error) {
	/*
	err := gomod.CheckPath(path)
	if err != nil {
		return fmt.Errorf("bad module name %q, should have domain format 'domain.com/...'", path)
	}
	*/

	if path == cm.Module {
		return fmt.Errorf("cannot get current module")
	}

	// check for indirect and delete, the user is making it direct
	_, ok := cm.Indirect[path]
	if ok {
		delete(cm.Indirect, path)
	}

	// if latest, update
	if ver == "latest" || ver == "next" {
		ver, err = cache.GetLatestTag(path, gflags.Prerelease || ver == "next")
		if err != nil {
			return err
		}
		fmt.Println("found:", ver)
	}

	// check for already required at a version equal or greater (no downgrades with get)
	// HMMM, let's remove this restriction, GO supports downgrading, let's inform the user
	currVer, ok := cm.Require[path]
	if ok {
		if semver.Compare(currVer, ver) > 0 {
			fmt.Printf("Downgrading: %s from %s to %s\n", path, currVer, ver)
		}
	}

	// add dep to required
	cm.Require[path] = ver


	return nil
}


func updateAll(cm *CueMod, rflags flags.RootPflagpole, gflags flags.Mod__GetFlagpole) (err error) {
	for path, ver := range cm.Require {
		nver, err := cache.GetLatestTag(path, false)
		if err != nil {
			return err
		}

		// only update if newer, incase we have specific prereleases
		if semver.Compare(nver, ver) == 1 {
			cm.Require[path] = nver
		}
	}

	for path, ver := range cm.Indirect {
		// explicitly not doing prerelease here, should we allow updating all to pre-releases?
		nver, err := cache.GetLatestTag(path, false)
		if err != nil {
			return err
		}

		// only update if newer, incase we have specific prereleases
		if semver.Compare(nver, ver) > 1 {
			cm.Indirect[path] = ver
		}
	}

	return nil
}
