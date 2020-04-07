package gen

import (
	"io/ioutil"
	"path"

	"github.com/hofstadter-io/hof/lib/util"
)

func (F *File) WriteOutput() error {
	var err error

	// fmt.Println("WriteFile:", F.Filepath)
	// fmt.Printf("%#+v\n\n", F)

	err = util.Mkdir(path.Join(F.Filepath))
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

	fn := path.Join(basedir, F.Filepath)

	err = util.Mkdir(fn)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fn, F.RenderContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

