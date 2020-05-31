package modder

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	gomod "golang.org/x/mod/module"

	"github.com/hofstadter-io/hof/lib/mod/parse/modfile"
	"github.com/hofstadter-io/hof/lib/yagu"
)

/* TODO
- more configuration for intialization
*/

func (mdr *Modder) Init(module string) error {
	// exec commands override configurable behavior
	if len(mdr.CommandInit) > 0 {
		for _, cmd := range mdr.CommandGraph {
			out, err := yagu.Exec(cmd)
			fmt.Println(out)
			if err != nil {
				return err
			}
		}
	} else {
		// or load from file
		err := mdr.initModFile(module)
		if err != nil {
			return err
		}
	}

	// Run init templates regardless of the init main method
	err := mdr.writeInitTemplates(module)
	if err != nil {
		return err
	}

	return nil
}

func (mdr *Modder) initModFile(module string) error {
	lang := mdr.Name
	filename := mdr.ModFile

	err := gomod.CheckPath(module)
	if err != nil {
		return fmt.Errorf("bad module format %q, should be 'domain.com/repo/proj'", module)
	}

	// make sure file does not exist
	_, err = ioutil.ReadFile(filename)
	// we read the file and it exists
	if err == nil {
		return fmt.Errorf("%s already exists", filename)
	}
	// error was not path error, so return
	if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
		return err
	}

	// Create empty modfile
	f, err := modfile.Parse(filename, nil, nil)
	if err != nil {
		return err
	}

	err = f.AddModuleStmt(module)
	if err != nil {
		return err
	}

	err = f.AddLanguageStmt(lang, mdr.Version)
	if err != nil {
		return err
	}

	bytes, err := f.Format()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (mdr *Modder) writeInitTemplates(module string) error {

	for filename, templateStr := range mdr.InitTemplates {

		tmpl, err := template.New(filename).Parse(templateStr)
		if err != nil {
			return err
		}

		data := map[string]interface{}{
			"Language": mdr.Name,
			"Module":   module,
			"Modder":   mdr,
		}

		err = os.MkdirAll(path.Dir(filename), 0755)
		if err != nil {
			return err
		}

		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		err = tmpl.Execute(file, data)
		if err != nil {
			return err
		}

	}

	return nil
}
