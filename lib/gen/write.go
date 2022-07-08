package gen

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func (F *File) WriteOutput() error {
	var err error

	if F.Filepath == "-" || strings.HasPrefix(F.Filepath, "hof-stdout-") {
		fmt.Print(string(F.FinalContent))
		return nil
	}

	dir := path.Dir(F.Filepath)
	err = yagu.Mkdir(dir)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(F.Filepath, F.FinalContent, 0644)
	if err != nil {
		return err
	}

	F.IsWritten = 1

	return nil
}

func (F *File) WriteShadow(basedir string) error {
	var err error

	// need to join and then find dir
	fn := path.Join(basedir, F.Filepath)
	dir := path.Dir(fn)

	err = yagu.Mkdir(dir)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fn, F.RenderContent, 0644)
	if err != nil {
		return err
	}

	return nil
}
