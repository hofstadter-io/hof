package mod

import (
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/repos/cache"
)

func Tidy(rflags flags.RootPflagpole, mflags flags.ModPflagpole) (error) {
	upgradeHofMods()

	cm, err := loadRootMod()
	if err != nil {
		return err
	}

	if len(cm.Require) == 0 {
		return nil
	}

	fns := []func () error {
		cm.SolveMVS,
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

	return cm.Vendor(true, rflags.Verbosity)
}

func (cm *CueMod) CleanDeps() error {
	// fmt.Println("clean deps:", cm.Module)

	// get ready for swap
	orig := cm.Require
	cm.Require  = make(map[string]string)
	cm.Indirect = make(map[string]string)

	// loop over build list, filling require & indirect
	for _, dep := range cm.BuildList {
		if dep.Path == cm.Module || strings.HasPrefix(dep.Path, ".") {
			continue
		}
		_, ok := orig[dep.Path]
		if ok {
			cm.Require[dep.Path] = dep.Version
		} else {
			cm.Indirect[dep.Path] = dep.Version
		}
	}

	// add any missing required back, as they were likely replaced
	for path, ver := range orig {
		if _, ok := cm.Require[path]; !ok {
			cm.Require[path] = ver
		}
	}

	return nil
}

func (cm *CueMod) CleanSums() error {
	// fmt.Println("clean sums:", cm.Module)
	keep := make(map[Dep][]string)

	// first, remove any sums we don't know about
	for dep, hashes := range cm.Sums {
		rver, rok := cm.Require[dep.Path]
		iver, iok := cm.Indirect[dep.Path]

		// found match in required
		if rok && dep.Version == rver {
			keep[dep] = hashes
			continue
		}

		// found match in indirect
		if iok && dep.Version == iver {
			keep[dep] = hashes
			continue
		}
	}
	
	// second, add any required, indirect, or replaces we have
	for path, ver := range cm.Require {
		dep := Dep{ Path: path, Version: ver }
		_, ok := keep[dep]
		if !ok {
			// fmt.Println("adding:", dep)
			h1, err := cache.Checksum(path, ver)
			if err != nil {
				return err
			}
			keep[dep] = []string{h1}
		}
	}
	for path, ver := range cm.Indirect {
		dep := Dep{ Path: path, Version: ver }
		_, ok := keep[dep]
		if !ok {
			// fmt.Println("adding:", dep)
			h1, err := cache.Checksum(path, ver)
			if err != nil {
				return err
			}
			keep[dep] = []string{h1}
		}
	}
	for _, repl := range cm.Replace {
		if repl.Version != "" && !strings.HasPrefix(repl.Path, ".") {
			dep := Dep{ Path: repl.Path, Version: repl.Version }
			_, ok := keep[dep]
			if !ok {
				// fmt.Println("adding:", dep)
				h1, err := cache.Checksum(repl.Path, repl.Version)
				if err != nil {
					return err
				}
				keep[dep] = []string{h1}
			}
		}
	}

	// finally update sums
	cm.Sums = keep

	return nil
}
