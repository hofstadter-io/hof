package gen

import (
	"fmt"
	"path/filepath"
	"os"
)

const SHADOW_DIR = ".hof"

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

	err = filepath.Walk(SHADOW_DIR, func(path string, info os.FileInfo, err error) error {
		// Don't need to save directories
		if info.IsDir() {
			return nil
		}

		if verbose {
			fmt.Println("  adding:", info.Name())
		}

		shadow[path] = &File {
			Filename: info.Name(),
		}

		return nil
	})

	if err != nil {
		err = fmt.Errorf("error walking the shadow dir %q: %w\n", SHADOW_DIR, err)
		return nil, err
	}

	return shadow, nil
}

