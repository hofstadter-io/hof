package yagu

import (
	"os"
	"path/filepath"
)

// cleanup empty dirs, walk up
func RemoveEmptyDirs(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	// if entries, we can return
	if len(entries) > 0 {
		return nil
	}

	// remove dir
	err = os.Remove(dir)
	if err != nil {
		return err
	}

	// recurse to parent dir
	return RemoveEmptyDirs(filepath.Dir(dir))
}
