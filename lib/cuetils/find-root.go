package cuetils

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindModuleAbsPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

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
		return "", fmt.Errorf("unable to find CUE module root")
	}

	return dir, nil
}
