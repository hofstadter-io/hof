package templates

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aymerick/raymond"
)

type Template struct {
	*raymond.Template
}

func NewTemplate() *Template {
	return &Template{}
}

func NewMap() TemplateMap {
	return NewTemplateMap()
}

func (template *Template) Render(design interface{}) (string, error) {
	tpl := (*raymond.Template)(template.Template)
	return tpl.Exec(design)
}

type TemplateMap map[string]*Template

func NewTemplateMap() TemplateMap {
	return make(map[string]*Template)
}


/*
Where's your docs doc?!
*/
func RenderTemplate(template *Template, design interface{}) (output string, err error) {
	tpl := (*raymond.Template)(template.Template)
	return tpl.Exec(design)
}

/*
Where's your docs doc?!
*/
func AddHelpersToRaymond(tpl *Template) {
	rtpl := (*raymond.Template)(tpl.Template)
	addTemplateHelpers(rtpl)
}

/*
Where's your docs doc?!
*/
func AddHelpersToTemplate(tpl *Template) {
	rtpl := (*raymond.Template)(tpl.Template)
	addTemplateHelpers(rtpl)
}

/*
Where's your docs doc?!
*/
func CreateTemplateFromString(name, source string) (tpl *Template, err error) {
	// parse template
	rtpl, err := raymond.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("While parsing: %s\n%w\n", name, err)
	}

	addTemplateHelpers(rtpl)

	return &Template{rtpl}, nil
}

func CreateTemplateFromFile(filename string) (tpl *Template, err error) {
	raw_template, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	source := string(raw_template)

	return CreateTemplateFromString(filename, source)
}

/*
Where's your docs doc?!
*/
func CreateTemplateMapFromFolder(folder string) (tplMap TemplateMap, err error) {
	tplMap = NewTemplateMap()
	err = tplMap.ImportFromFolder(folder)
	if err != nil {
		return nil, fmt.Errorf("while importing %s\n%w\n", folder, err)
	}
	return tplMap, nil
}

// HOFSTADTER_BELOW

func (M TemplateMap) ImportTemplateFile(filename string) error {
	return M.import_template("", filename)
}

func (M TemplateMap) ImportFromFolder(folder string) error {
	import_template_walk_func := func(base_path string) filepath.WalkFunc {
		return func(path string, info os.FileInfo, err error) error {
			local_m := M
			if err != nil {
				if _, ok := err.(*os.PathError); ok {
					return nil
				}
				return err
			}
			if info.IsDir() {
				return nil
			}

			return local_m.import_template(base_path, path)
		}
	}

	// Walk the directory
	err := filepath.Walk(folder, import_template_walk_func(folder))
	if err != nil {
		return err
	}
	return nil
}

func (M TemplateMap) import_template(base_path, path string) error {
	tpl_name := path
	L := len(base_path)
	if L > 0 {
		// should handle trailing slashes better here
		if base_path[L-1] == '/' {
			tpl_name = path[L:]
		} else {
			tpl_name = path[L+1:]
		}
	}

	raw_template, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	source := string(raw_template)

	// parse template
	tpl, err := raymond.Parse(source)
	if err != nil {
		return fmt.Errorf("While parsing file: %s\n%w", tpl_name, err)
	}

	addTemplateHelpers(tpl)

	M[tpl_name] = &Template{tpl}
	return nil
}
