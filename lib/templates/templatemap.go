package templates

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
)

type TemplateMap map[string]*Template

func NewTemplateMap() TemplateMap {
	return make(map[string]*Template)
}

func CreateTemplateMapFromFolder(glob, prefix string, delims Delims, delimMap map[string]Delims) (tplMap TemplateMap, err error) {
	tplMap = NewTemplateMap()
	err = tplMap.ImportFromFolder(glob, prefix, delims, delimMap)
	if err != nil {
		return nil, fmt.Errorf("while importing %s\n%w\n", glob, err)
	}
	return tplMap, nil
}

func (M TemplateMap) ImportTemplateFile(filename string, delims Delims) error {
	return M.importTemplate(filename, "", delims)
}

func (M TemplateMap) ImportFromFolder(glob, prefix string, delims Delims, delimMap map[string]Delims) error {
	// all templates
	matches, err := zglob.Glob(glob)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		// return fmt.Errorf("No templates found for '%s'", glob)
		return nil
	}

	// delimOverrides
	overrides := make(map[string]Delims)
	for g, d := range delimMap {
		gmatches, err := zglob.Glob(g)
		if err != nil {
			return err
		}
		if len(gmatches) > 0 {
			for _, fn := range gmatches {
				overrides[fn] = d
			}
		}
	}

	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			return err
		}

		if info.IsDir() {
			continue
		}

		// override delims?
		d := delims
		if len(overrides) > 0 {
			D, ok := overrides[match]
			if ok {
				d = D
			}
		}

		err = M.importTemplate(match, prefix, d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (M TemplateMap) importTemplate(filePath, prefix string, delims Delims) error {
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	content := string(source)

	T, err := CreateFromString(filePath, content, delims)
	if err != nil {
		return fmt.Errorf("While parsing template file: %s\n%w", filePath, err)
	}

	// clean up filename before inserting into map, so when users reference in their generators, we align
	filePath = strings.TrimPrefix(filePath, prefix)
	filePath = strings.TrimPrefix(filePath, "/")  // be resilent to trailing slashes in the last value, or the lack therein
	filePath = filepath.Clean(filePath)

	M[filePath] = T
	return nil
}
