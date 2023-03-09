package mod

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue/parser"

	"github.com/hofstadter-io/hof/lib/repos/utils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

var wellKnownDirs = []string {
	".hof/",
	".workdir/",
	"cue.mod/",
	"docs/",
	"testdata/",
}

func (cm *CueMod) findDepsFromImports() (err error) {

	// dirs, err := cm.findDirsWithCueFiles()
	files, err := cm.findCueFiles()
	if err != nil {
		return err
	}

	// fmt.Println("finding missing:")
	depMap := make(map[string]bool)
	// for _, dir := range dirs {
		// deps, _ := cm.dir2deps(dir)
	for _, file := range files {
		// fmt.Println("dir:", dir)
		deps, _ := cm.file2deps(file)
		for _, dep := range deps {
			// check to see if we know about dep already
			if cm.checkIfDepAlreadyKnown(dep) {
				continue
			}

			// memo if we haven't already
			if _, ok := depMap[dep]; !ok {
				depMap[dep] = true
			}
		}
	}

	if len(depMap) == 0 {
		return nil
	}

	// fmt.Println("adding deps: ")
	for dep, _ := range depMap {
		// fmt.Println("  d:", dep)
		// need to walk up to find
		m, err := cm.dep2module(dep)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// fmt.Println("  ", m)
		cm.Require[m] = "latest"
	}

	return nil
}

func (cm *CueMod) checkIfDepAlreadyKnown(dep string) (found bool) {
	// from this module
	if dep == cm.Module || strings.HasPrefix(dep, cm.Module + "/") {
		return true
	}

	for path, _ := range cm.Replace {
		if dep == path || strings.HasPrefix(dep, path + "/") {
			return true
		}
	}

	for path, _ := range cm.Require {
		if dep == path || strings.HasPrefix(dep, path + "/") {
			return true
		}
	}

	for path, _ := range cm.Indirect {
		if dep == path || strings.HasPrefix(dep, path + "/") {
			return true
		}
	}

	return false
}

func (cm *CueMod) findCueFiles() (files []string, err error) {
	// fmt.Println("findDirs:", cm.Basedir)

	files, err = yagu.FilesFromGlobs([]string{filepath.Join(cm.Basedir, "**/*.cue")})
	if err != nil {
		return nil, err
	}

	final := make([]string, 0, len(files))

	for _, file := range files {
		// skip wellknown dirs
		match := false
		for _, wkd := range wellKnownDirs {
			if strings.Contains(file, wkd) {
				match = true
				break
			}
		}
		if match {
			continue
		}

		final = append(final, file)
	}

	return final, nil
}

func (cm *CueMod) file2deps(file string) ([]string, []error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, []error{err}
	}

	fn, err := filepath.Rel(cm.Basedir, file)
	if err != nil {
		return nil, []error{err}
	}

	f, err := parser.ParseFile(fn, content)
	if err != nil {
		return nil, []error{err}
	}

	paths := []string{}

	for _, imp := range f.Imports {
		// this is the import path string value
		path := imp.Path.Value

		// remove the quotes
		path = strings.Replace(path, "\"", "", -1)

		// filter builtin packages
		if !strings.Contains(path, ".") {
			continue
		}

		// include in our deps
		paths = append(paths, path)
	}

	return paths, nil
}

func (cm *CueMod) dep2module(dep string) (string, error) {
	remote, owner, repo := utils.ParseModURL(dep)

	m := path.Join(remote, owner, repo)
	return m, nil
}
