package util

import (
	"archive/tar"
	"compress/gzip"
	// "fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
)

func UntarFiles(files []string, dst string, content []byte) error {

	if _, err := os.Stat(tarfile); !os.IsNotExist(err) {
		// path/to/whatever does not exist
		err = os.RemoveAll(tarfile)
		if err != nil {
			return err
		}
	}


	for _, dir := range files {
		err := os.RemoveAll(filepath.Join(dst, dir))
		if err != nil {
			return err
		}
	}

  err := ioutil.WriteFile(tarfile, content, 0644)
	if err != nil {
		return err
	}

	err = archiver.Unarchive(tarfile, ".")
	if err != nil {
		return err
	}

	/*
	err = os.RemoveAll(tarfile)
	if err != nil {
		return data, err
	}
	*/

	return nil
}

func OldUntarFiles(files []string, dst string, r io.Reader) error {

	for _, dir := range files {
		err := os.RemoveAll(filepath.Join(dst, dir))
		if err != nil {
			return err
		}
	}


	gzr, err := gzip.NewReader(r)
	if err != nil {

	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			dir := filepath.Dir(target)
			// fmt.Println("x ", target)
			if _, err := os.Stat(dir); err != nil {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return err
				}
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}


