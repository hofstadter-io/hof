package templates

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mattn/go-zglob"
)

type TemplateMap map[string]*Template

func NewTemplateMap() TemplateMap {
	return make(map[string]*Template)
}

func CreateTemplateMapFromFolder(glob, prefix string, delims *Delims) (tplMap TemplateMap, err error) {
	tplMap = NewTemplateMap()
	err = tplMap.ImportFromFolder(glob, prefix, delims)
	if err != nil {
		return nil, fmt.Errorf("while importing %s\n%w\n", glob, err)
	}
	return tplMap, nil
}

func (M TemplateMap) ImportTemplateFile(filename string, delims *Delims) error {
	return M.importTemplate(filename, "", delims)
}

func (M TemplateMap) ImportFromFolder(glob, prefix string, delims *Delims) error {
	matches, err := zglob.Glob(glob)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return fmt.Errorf("No templates found for '%s'", glob)
	}

	for _, match := range matches {
		err := M.importTemplate(match, prefix, delims)
		if err != nil {
			return err
		}
	}

	return nil
}

func (M TemplateMap) importTemplate(filePath, prefix string, delims *Delims) error {
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	content := string(source)

	T, err := CreateFromString(filePath, content, delims)
	if err != nil {
		return fmt.Errorf("While parsing template file: %s\n%w", filePath, err)
	}

	if prefix != "" {
		filePath, _ = filepath.Rel(prefix, filePath)
	} else {
		filePath = filepath.Clean(filePath)
	}
	M[filePath] = T
	return nil
}
