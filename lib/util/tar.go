package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver"
)

var tarfile = "studios.tar.gz"

func TarBallFiles(files []string, src, tarFile string) (err error) {

	if _, err := os.Stat(tarFile); !os.IsNotExist(err) {
		// path/to/whatever does not exist
		err = os.RemoveAll(tarFile)
		if err != nil {
			return err
		}
	}

	var paths []string
	// ensure the src actually exists before trying to tar it
	for _, path := range files {
		lpath := filepath.Join(src, path)
		if _, lerr := os.Stat(lpath); lerr == nil {
			paths = append(paths, lpath)
		}
	}

	err = archiver.Archive(paths, tarFile)
	if err != nil {
		return err
	}

	return nil
}

func TarFiles(files []string, src string) (data []byte, err error) {

	err = TarBallFiles(files, src, tarfile)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadFile(tarfile)
	if err != nil {
		return nil, err
	}

	/*
		err = os.RemoveAll(tarfile)
		if err != nil {
			return data, err
		}
	*/

	return data, err
}

func OldTarFiles(files []string, src string, writers ...io.Writer) error {
	var paths []string
	// ensure the src actually exists before trying to tar it
	for _, path := range files {
		if _, err := os.Stat(filepath.Join(src, path)); err == nil {
			paths = append(paths, path)
		}
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk path
	walker := func(file string, fi os.FileInfo, err error) error {
		// fmt.Println("+", file)

		// return on any error
		if err != nil {
			return err
		}

		// return on non-regular files (thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update)
		if !fi.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	}

	for _, path := range paths {
		err := filepath.Walk(filepath.Join(src, path), walker)
		if err != nil {
			fmt.Println("Error", path, err)
		}
	}

	return nil
}
