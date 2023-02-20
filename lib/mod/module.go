package mod

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/go-git/go-billy/v5"
	"golang.org/x/mod/module"
)



type CueMod struct {
	Basedir string

	Module string

	CueVer string

	Require  map[string]string
	Indirect map[string]string
	Replace  map[string]Dep
	Sums     map[Dep][]string

	// final list produced by MVS
	BuildList []module.Version

	FS       billy.Filesystem
}

type Dep struct {
	Path    string
	Version string
}

func ReadModule(basedir string, FS billy.Filesystem) (cm *CueMod, err error) {
	cm = new(CueMod)

	cm.Basedir = basedir
	cm.FS = FS

	err = cm.ReadModFile()
	if err != nil {
		return cm, err
	}

	err = cm.ReadSumFile()
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && (strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file")) {
			return cm, err
		}
	}

	return cm, nil
}

func (cm *CueMod) ReadModFile() (err error) {
	var data []byte
	fn := filepath.Join("cue.mod", "module.cue")

	f, err := cm.FS.Open(fn)
	if err != nil {
		return err
	}

	data, err = io.ReadAll(f)
	if err != nil {
		return err
	}
	return cm.ParseModFile(data)
}

func (cm *CueMod) ParseModFile(data []byte) (err error) {

	c := cuecontext.New()

	v := c.CompileBytes(data)

	// TODO, consider looping over all fields and using a switch
	// it's how we could extract extra metadata for any tool

	// read metadata files
	cm.Module, err = v.LookupPath(cue.ParsePath("module")).String()
	if err != nil {
		return err
	}
	cm.CueVer, err = v.LookupPath(cue.ParsePath("cue")).String()
	if err != nil {
		return err
	}


	// read require files
	requires  := v.LookupPath(cue.ParsePath("require"))
	if requires.Exists() {
		err = requires.Decode(&(cm.Require))
		if err != nil {
			return err
		}
	}

	// parse indirects
	indirects := v.LookupPath(cue.ParsePath("indirect"))
	if indirects.Exists() {
		err = indirects.Decode(&(cm.Indirect))
		if err != nil {
			return err
		}
	}

	// parse replaces
	cm.Replace = make(map[string]Dep)
	replaces := v.LookupPath(cue.ParsePath("replace"))
	if replaces.Exists() {
		iter, err := replaces.Fields()
		if err != nil {
			return err
		}
		for iter.Next() {
			label := iter.Label()
			value := iter.Value()

			switch value.IncompleteKind() {

			case cue.StringKind:
				info, err := value.String()
				if err != nil {
					return err
				}
				dep := Dep{ Path: info }
				cm.Replace[label] = dep

			case cue.StructKind:
				info := make(map[string]string)
				err := value.Decode(&info)
				if err != nil {
					return err
				}

				for k, v := range info {
					rep := Dep{ Path: k, Version: v }
					cm.Replace[label] = rep
				}
			}

		}
	}

	return nil
}

func (cm *CueMod) ReadSumFile() error {
	fn := filepath.Join(cm.Basedir, "cue.mod", "sums.cue")

	data, err := os.ReadFile(fn)
	if err != nil {
		return err
	}

	c := cuecontext.New()

	v := c.CompileBytes(data)

	sums := v.LookupPath(cue.ParsePath("sums"))

	cm.Sums = make(map[Dep][]string)
	iter, err := sums.Fields()
	if err != nil {
		return err
	}

	for iter.Next() {
		label := iter.Label()
		value := iter.Value()

		info := make(map[string][]string)
		err := value.Decode(&info)
		if err != nil {
			return err
		}

		for k, v := range info {
			lhs := Dep{ Path: label, Version: k }
			cm.Sums[lhs] = v
		}
	}

	return nil
}

func (cm *CueMod) WriteModule() (err error) {
	err = cm.WriteModFile()
	if err != nil {
		return err
	}
	return cm.WriteSumFile()
}

func (cm *CueMod) WriteModFile() (err error) {

	var buf bytes.Buffer

	// write top-level data
	buf.WriteString(fmt.Sprintf("module: %q\n", cm.Module))
	buf.WriteString(fmt.Sprintf("cue: %q\n", cm.CueVer))


	// write requires
	var sorted []module.Version
	for path, ver := range cm.Require {
		sorted = append(sorted, module.Version{Path: path, Version: ver})
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Path == sorted[j].Path {
			return sorted[i].Version < sorted[j].Version
		}
		return sorted[i].Path < sorted[j].Path
	})
	buf.WriteString("\nrequire: {\n")
	for _, ver := range sorted {
		m := fmt.Sprintf("\t%q: %q\n", ver.Path, ver.Version)
		buf.WriteString(m)
	}
	buf.WriteString("}\n")

	// write indirects
	if len(cm.Indirect) > 0 {
		var sorted []module.Version
		for path, ver := range cm.Indirect {
			sorted = append(sorted, module.Version{Path: path, Version: ver})
		}
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].Path == sorted[j].Path {
				return sorted[i].Version < sorted[j].Version
			}
			return sorted[i].Path < sorted[j].Path
		})
		buf.WriteString("\nindirect: {\n")
		for _, ver := range sorted {
			m := fmt.Sprintf("\t%q: %q\n", ver.Path, ver.Version)
			buf.WriteString(m)
		}
		buf.WriteString("}\n")
	}

	// write replaces
	if len(cm.Replace) > 0 {
		var sorted []module.Version
		for path, _ := range cm.Replace {
			sorted = append(sorted, module.Version{Path: path})
		}
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Path < sorted[j].Path
		})
		buf.WriteString("\nreplace: {\n")
		for _, path := range sorted {
			dep := cm.Replace[path.Path]
			m := fmt.Sprintf("\t%q: %q: %q\n", path.Path, dep.Path, dep.Version)
			if dep.Version == "" {
				m = fmt.Sprintf("\t%q: %q\n", path.Path, dep.Path)
			}
			buf.WriteString(m)
		}
		buf.WriteString("}\n")
	}

	return os.WriteFile(filepath.Join("cue.mod/module.cue"), buf.Bytes(), 0644)
}

func (cm *CueMod) WriteSumFile() (err error) {
	// build up slice
	var sorted []Dep
	for ver, _ := range cm.Sums {
		sorted = append(sorted, ver)
	}

	// sort slice by ver.Path
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Path == sorted[j].Path {
			return sorted[i].Version < sorted[j].Version
		}
		return sorted[i].Path < sorted[j].Path
	})

	var buf bytes.Buffer
	buf.WriteString("sums: {\n")

	for _, ver := range sorted {
		hashes := cm.Sums[ver]
		m := fmt.Sprintf("\t%q: %q: ", ver.Path, ver.Version)
		h, err := json.Marshal(hashes)
		if err != nil {
			return err
		}

		buf.WriteString(m)
		buf.Write(h)
		buf.WriteRune('\n')
	}

	buf.WriteString("}\n")

	return os.WriteFile(filepath.Join("cue.mod/sums.cue"), buf.Bytes(), 0644)
}
