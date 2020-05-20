package yagu

import (
	"fmt"
	"os"
)

func Mkdir(dir string) error {
	var err error

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

