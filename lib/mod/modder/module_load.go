package modder

import (
	"fmt"
	"os"

	"github.com/hofstadter-io/hof/lib/mod/parse/mappingfile"
	"github.com/hofstadter-io/hof/lib/mod/parse/modfile"
	"github.com/hofstadter-io/hof/lib/mod/parse/sumfile"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func (m *Module) LoadModFile(fn string, ignoreReplace bool) error {

	modBytes, err := yagu.BillyReadAll(fn, m.FS)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
			return err
		}
	} else {
		f, err := modfile.Parse(fn, modBytes, nil)
		if err != nil {
			return err
		}
		m.ModFile = f
		m.Language = f.Language.Name
		m.LangVer = f.Language.Version
		// m.Module = f.Module.Mod.Path
		// m.Version = f.Module.Mod.Version

		for _, req := range f.Require {
			r := Require{Path: req.Mod.Path, Version: req.Mod.Version}
			m.Require = append(m.Require, r)
		}

		// Let's just not load them for now to be sure
		if !ignoreReplace {
			for _, rep := range f.Replace {
				r := Replace{OldPath: rep.Old.Path, OldVersion: rep.Old.Version, NewPath: rep.New.Path, NewVersion: rep.New.Version}
				m.Replace = append(m.Replace, r)
			}
		}

		err = m.MergeSelfDeps(ignoreReplace)
		if err != nil {
			return err
		}

	}

	return nil
}

func (m *Module) MergeSelfDeps(ignoreReplace bool) error {
	// Now merge self deps
	m.SelfDeps = map[string]Replace{}
	for _, req := range m.Require {
		if _, ok := m.SelfDeps[req.Path]; ok {
			return fmt.Errorf("Dependency %q required twice in %q", req.Path, m.Module)
		}
		m.SelfDeps[req.Path] = Replace{
			NewPath:    req.Path,
			NewVersion: req.Version,
		}
	}

	// we typically ignore replaces from when not the root module
	if !ignoreReplace {
		dblReplace := map[string]Replace{}
		for _, rep := range m.Replace {
			// Check if replaced twice
			if _, ok := dblReplace[rep.OldPath]; ok {
				return fmt.Errorf("Dependency %q replaced twice in %q", rep.OldPath, m.Module)
			}
			dblReplace[rep.OldPath] = rep

			// Pull in require info if not in replace
			if req, ok := m.SelfDeps[rep.OldPath]; ok {
				if rep.OldVersion == "" {
					rep.OldVersion = req.NewVersion
				}
			}
			m.SelfDeps[rep.OldPath] = rep
		}
	}

	return nil
}

func (m *Module) LoadSumFile(fn string) error {

	sumBytes, err := yagu.BillyReadAll(fn, m.FS)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
			return err
		}
	} else {
		sumMod, err := sumfile.ParseSum(sumBytes, fn)
		if err != nil {
			return err
		}
		m.SumFile = &sumMod
	}

	return nil
}

func (m *Module) LoadMappingFile(fn string) error {

	mapBytes, err := yagu.BillyReadAll(fn, m.FS)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
			return err
		}
	} else {
		mapMod, err := mappingfile.ParseMapping(mapBytes, fn)
		if err != nil {
			return err
		}
		m.Mappings = &mapMod
	}

	return nil
}

func (m *Module) LoadMetaFiles(modname, sumname, mapname string, ignoreReplace bool) error {
	var err error

	// TODO load the modules .mvsconfig if present

	err = m.LoadModFile(modname, ignoreReplace)
	if err != nil {
		return err
	}

	// Suppressing the next to errors here
	//   because we can handle them not being around else where
	err = m.LoadSumFile(sumname)
	if err != nil {
		// fmt.Println(err)
		// return err
	}

	err = m.LoadMappingFile(mapname)
	if err != nil {
		// fmt.Println(err)
		// return err
	}

	return nil
}
