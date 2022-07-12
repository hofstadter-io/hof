package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const SHADOW_DIR = ".hof/shadow/"

func LoadShadow(subdir string, verbosity int) (map[string]*File, error) {
	shadowDir := filepath.Join(SHADOW_DIR, subdir)
	shadow := map[string]*File{}

	if verbosity > 1 {
		fmt.Printf("Loading shadow @ %q\n", shadowDir)
	}

	_, err := os.Lstat(shadowDir)
	if err != nil {
		// make sure we check err for something actually bad
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" {
			return nil, err
		}
		// file not found, leave politely
		if verbosity > 1 {
			fmt.Println("  shadow not found")
		}
		return shadow, nil
	}

	err = filepath.Walk(shadowDir, func(fpath string, info os.FileInfo, err error) error {
		// Don't need to save directories
		if info.IsDir() {
			return nil
		}

		// read contents
		bytes, err := ioutil.ReadFile(fpath)
		if err != nil {
			return err
		}

		// simplify path
		fpath = strings.TrimPrefix(fpath, shadowDir + "/")

		// debug
		if verbosity > 1 {
			fmt.Println("  adding:", fpath)
		}

		shadow[fpath] = &File{
			FinalContent: bytes,
			Filepath: fpath,
		}

		return nil
	})

	if err != nil {
		err = fmt.Errorf("error walking the shadow dir %q: %w\n", shadowDir, err)
		return nil, err
	}

	return shadow, nil
}
