package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const SHADOW_DIR = ".hof/"

func LoadShadow(verbose bool) (map[string]*File, error) {
	if verbose {
		fmt.Printf("Loading shadow @ %q\n", SHADOW_DIR)
	}

	shadow := map[string]*File{}

	_, err := os.Lstat(SHADOW_DIR)
	if err != nil {
		// make sure we check err for something actually bad
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" {
			return nil, err
		}
		// file not found, leave politely
		if verbose {
			fmt.Println("  shadow not found")
		}
		return shadow, nil
	}

	err = filepath.Walk(SHADOW_DIR, func(fpath string, info os.FileInfo, err error) error {
		// Don't need to save directories
		if info.IsDir() {
			return nil
		}

		if verbose {
			fmt.Println("  adding:", info.Name())
		}

		fpath = strings.TrimPrefix(fpath, SHADOW_DIR)
		shadow[fpath] = &File {
			Filepath: fpath,
		}

		return nil
	})

	if err != nil {
		err = fmt.Errorf("error walking the shadow dir %q: %w\n", SHADOW_DIR, err)
		return nil, err
	}

	return shadow, nil
}

func (F *File) ReadShadow() error {
	if F.ShadowFile == nil {
		return nil
	}

	// Should have already been confirmed to exist at this point

	bytes, err := ioutil.ReadFile(path.Join(SHADOW_DIR, F.ShadowFile.Filepath))
	if err != nil {
		return err
	}

	F.ShadowFile.FinalContent = bytes

	return nil
}

