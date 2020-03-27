package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func (F *File) WriteOutput() error {
	var err error

	err = mkdir(path.Join(F.Filename))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(F.Filename, F.FinalContent, 0644)
	if err != nil {
		return err
	}

	F.IsWritten = 1

	return nil
}

func (F *File) WriteShadow() error {
	var err error

	fn := path.Join(SHADOW_DIR, F.Filename)

	err = mkdir(fn)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fn, F.RenderContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

func mkdir(filename string) error {
	var err error

	// get directory from filename
	dir := path.Dir(filename)

	// Let's look for the directory
	info, err := os.Lstat(dir)
	if err != nil {

		// make sure we check err for something actually bad
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" {
			return err
		}

		// file not found, good!
		// go to the mkdir below

	} else {

		// Something is there
		if info.IsDir() {
			// Our directory already exists
			return nil
		} else {
			// That something else
			return fmt.Errorf("Mkdir: %q exists but is not a directory", info.Name())
		}

	}

	// Make the directory
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	return nil
}
