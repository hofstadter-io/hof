package gen

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func (F *File) WriteOutput(basedir string) error {
	// print to stdout
	if F.Filepath == "-" || strings.HasPrefix(F.Filepath, "hof-stdout-") {
		fmt.Print(string(F.FinalContent))
		return nil
	}

	// write to file
	err := F.write(basedir, F.FinalContent)
	if err != nil {
		return err
	}

	F.IsWritten = 1

	return nil
}

func (F *File) WriteShadow(basedir string) error {
	return F.write(basedir, F.RenderContent)
}

func (F *File) write(basedir string, content []byte) error {

	fn := path.Join(basedir, F.Filepath)
	dir := path.Dir(fn)

	err := yagu.Mkdir(dir)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fn, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
