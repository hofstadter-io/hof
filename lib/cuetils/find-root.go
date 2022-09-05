package cuetils

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindModuleAbsPath(dir string) (string, error) {
	var err error
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	fmt.Println("finding root for:", dir)

	found := false

	for !found && dir != "/" {
		try := filepath.Join(dir, "cue.mod")
		info, err := os.Stat(try)
		if err == nil && info.IsDir() {
			found = true
			break
		}

		next := filepath.Clean(filepath.Join(dir, ".."))
		dir = next
	}

	if !found {
		return "", nil
		// return "", fmt.Errorf("unable to find CUE module root")
	}

	return dir, nil
}
