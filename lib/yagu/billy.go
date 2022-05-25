package yagu

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-billy/v5"
	"golang.org/x/mod/sumdb/dirhash"
)

func BillyReadAllString(filename string, FS billy.Filesystem) (string, error) {
	bytes, err := BillyReadAll(filename, FS)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func BillyReadAll(filename string, FS billy.Filesystem) ([]byte, error) {
	f, err := FS.Open(filename)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func BillyGetFilelist(FS billy.Filesystem) ([]os.FileInfo, error) {
	return FS.ReadDir("/")
}

// Writes dir in FS onto the os filesystem at baseDir
func BillyWriteDirToOS(baseDir string, dir string, FS billy.Filesystem) error {
	// fmt.Println("DIR:  ", baseDir, dir)
	files, err := FS.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		longname := path.Join(dir, file.Name())
		// fmt.Println("DIR:  ", baseDir, dir, file.Name(), longname)

		if file.IsDir() {
			err = BillyWriteDirToOS(baseDir, longname, FS)
			if err != nil {
				return err
			}

		} else {
			err = BillyWriteFileToOS(baseDir, longname, FS)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

// Writes file in FS onto the os filesystem at baseDir
func BillyWriteFileToOS(baseDir string, file string, FS billy.Filesystem) error {
	outName := path.Join(baseDir, file)

	// fmt.Println("FILE: ", outName)
	err := os.MkdirAll(path.Dir(outName), 0755)
	if err != nil {
		return err
	}

	src, err := FS.Open(file)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer dst.Close()

	io.Copy(dst, src)

	return nil
}

// Write dir in FS onto the os filesystem at baseDir
//
func BillyGlobWriteDirToOS(baseDir string, dir string, FS billy.Filesystem, includes, excludes []string) error {
	// fmt.Println("DIR:  ", baseDir, dir)
	files, err := FS.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		longname := path.Join(dir, file.Name())
		// fmt.Println("DIR:  ", baseDir, dir, file.Name(), longname, outname)
		// fmt.Println("GLOB?  ", longname)

		if file.IsDir() {
			err = BillyGlobWriteDirToOS(baseDir, longname, FS, includes, excludes)
			if err != nil {
				return err
			}

		} else {

			include, _ := CheckShouldInclude(longname, includes, excludes)
			// fmt.Println("COPY ==>", longname, include, exclude, include && !exclude)

			if include {
				err = BillyWriteFileToOS(baseDir, longname, FS)
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func BillyLoadFromZip(zReader *zip.Reader, FS billy.Filesystem, trimFirstDir bool) error {
	for _, f := range zReader.File {

		// Is this a directory?
		if strings.HasSuffix(f.Name, "/") {
			err := FS.MkdirAll(f.Name, 0755)
			if err != nil {
				return err
			}
			continue
		}

		fn := f.Name[strings.Index(f.Name, "/")+1:]

		src, err := f.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := FS.Create(fn)
		if err != nil {
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}

		// fmt.Println("FromZip:", fn, cnt)
	}

	return nil
}

// Loads an initialized  zip.Reader into an initialize billy.Filesystem
func BillyGlobLoadFromZip(zReader *zip.Reader, FS billy.Filesystem, includes, excludes []string) error {

	for _, f := range zReader.File {
		include, _ := CheckShouldInclude(f.Name, includes, excludes)
		if !include {
			continue
		}

		// Is this a directory?
		if strings.HasSuffix(f.Name, "/") {
			err := FS.MkdirAll(f.Name, 0755)
			if err != nil {
				return err
			}
			continue
		}

		fn := f.Name[strings.Index(f.Name, "/")+1:]

		src, err := f.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := FS.Create(fn)
		if err != nil {
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}

	}

	return nil
}

func BillyFilenames(dir string, FS billy.Filesystem) ([]string, error) {
	out := []string{}

	files, err := FS.ReadDir(dir)
	if err != nil {
		return out, err
	}

	for _, file := range files {
		fullname := path.Join(dir, file.Name())

		if file.IsDir() {
			tmp, err := BillyFilenames(fullname, FS)
			if err != nil {
				return out, err
			}
			out = append(out, tmp...)
		} else {
			out = append(out, fullname)
		}
	}

	return out, nil
}

func BillyCalcHash(FS billy.Filesystem) (string, error) {
	return BillyCalcDirHash("/", FS)
}

func BillyGlobCalcHash(FS billy.Filesystem, include, exclude []string) (string, error) {
	return BillyGlobCalcDirHash("/", FS, include, exclude)
}

func BillyCalcDirHash(dir string, FS billy.Filesystem) (string, error) {
	all, err := BillyFilenames(dir, FS)
	if err != nil {
		return "", err
	}

	var files []string
	for _, f := range all {
		// fmt.Println("FILE:", f)
		if strings.HasPrefix(f, "/.git/") || strings.HasPrefix(f, ".git/") {
			continue
		}

		files = append(files, f)
	}

	open := func(fn string) (io.ReadCloser, error) {
		return FS.Open(fn)
	}

	return dirhash.Hash1(files, open)
}

func BillyGlobCalcDirHash(dir string, FS billy.Filesystem, includes, excludes []string) (string, error) {
	all, err := BillyFilenames(dir, FS)
	if err != nil {
		return "", err
	}

	var files []string
	for _, f := range all {
		include, _ := CheckShouldInclude(f, includes, excludes)
		if !include {
			continue
		}
		// fmt.Println("FILE:", f)
		files = append(files, f)
	}

	open := func(fn string) (io.ReadCloser, error) {
		return FS.Open(fn)
	}

	return dirhash.Hash1(files, open)
}

func BillyCalcFileHash(filename string, FS billy.Filesystem) (string, error) {
	if !strings.HasPrefix(filename, "/") {
		filename = "/" + filename
	}

	open := func(fn string) (io.ReadCloser, error) {
		return FS.Open(fn)
	}

	return dirhash.Hash1([]string{filename}, open)
}
