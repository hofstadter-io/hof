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

func CreateTemplateMapFromFolder(folder, system string, config *Config, configGlobs map[string]*Config) (tplMap TemplateMap, err error) {
	tplMap = NewTemplateMap()
	err = tplMap.ImportFromFolder(folder, system, config, configGlobs)
	if err != nil {
		return nil, fmt.Errorf("while importing %s\n%w\n", folder, err)
	}
	return tplMap, nil
}

func (M TemplateMap) ImportTemplateFile(filename, system string, config *Config) error {
	return M.import_template("", filename, system, config)
}

func (M TemplateMap) ImportFromFolder(folder, system string, config *Config, configGlobs map[string]*Config) error {
	import_template_walk_func := func(base_path string) filepath.WalkFunc {
		return func(path string, info os.FileInfo, err error) error {
			// fmt.Println("templates.ImportFromFolder", path)
			local_m := M
			if err != nil {
				if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" {
					return err
				}
				return nil
			}
			if info.IsDir() {
				return nil
			}

			c, err := LookupConfig(path, configGlobs)
			if err != nil {
				return err
			}
			if c != nil {
				config = c
			}
			return local_m.import_template(base_path, path, system, config)
		}
	}

	// Walk the directory
	err := filepath.Walk(folder, import_template_walk_func(folder))
	if err != nil {
		return err
	}
	return nil
}

func (M TemplateMap) import_template(basePath, filePath, system string, config *Config) error {
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	content := string(source)

	// create template
	// TODO, how to explicitly set template system
	T, err := CreateFromString(filePath, content, system, config)
	if err != nil {
		return fmt.Errorf("While parsing template file: %s\n%w", filePath, err)
	}

	relFilePath := strings.TrimPrefix(filePath, basePath)
	M[relFilePath] = T
	return nil
}

func LookupConfig(fn string, cfgs map[string]*Config) (*Config, error) {
	fmt.Println("Lookup", fn)
	for glob, cfg := range cfgs {
		match, err := zglob.Match(glob, fn)
		if err != nil {
			return nil, err
		}

		if match {
			return cfg, nil
		}
	}
	return nil, nil
}
