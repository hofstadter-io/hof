package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
		fpath = strings.TrimPrefix(fpath, shadowDir + "/")

		if verbosity > 1 {
			fmt.Println("  adding:", fpath)
		}

		shadow[fpath] = &File{
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

func (F *File) ReadShadow() error {
	if F.ShadowFile == nil {
		return nil
	}

	// Should have already been confirmed to exist at this point
	shadowFN := path.Join(SHADOW_DIR, F.ShadowFile.Filepath)
	// fmt.Println("ReadShadow", shadowFN)
	bytes, err := ioutil.ReadFile(shadowFN)
	if err != nil {
		return err
	}

	F.ShadowFile.FinalContent = bytes

	return nil
}
