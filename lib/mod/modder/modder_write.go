package modder

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/hofstadter-io/yagu"

	"github.com/hofstadter-io/hof/lib/mod/parse/sumfile"
	"github.com/hofstadter-io/hof/lib/mod/util"
)

var (
	// Common files to copy from modules, also includes the .md version of the filename
	definiteVendors = []string{
		"README",
		"SECURITY",
		"AUTHORS",
		"CONTRIBUTORS",
		"COPYLEFT",
		"COPYING",
		"COPYRIGHT",
		"LEGAL",
		"LICENSE",
		"NOTICE",
		"PATENTS",
	}
)

func (mdr *Modder) WriteVendor() error {
	// fmt.Println("Writing Vendor from scratch")
	os.RemoveAll(mdr.ModsDir)

	// make vendor dir if not present
	err := yagu.Mkdir(mdr.ModsDir)
	if err != nil {
		return err
	}

	// write out each dep
	for _, m := range mdr.depsMap {
		// XXX, this only (?) happens with local replaces with no matching require entry
		if m.Version == "" {
			m.Version = "v0.0.0"
		}

		dirhash, err := util.BillyCalcHash(m.FS)
		if err != nil {
			mdr.errors = append(mdr.errors, err)
			return fmt.Errorf("While calculating billy dir hash for %q\n%w\n", mdr.ModsDir, err)
		}

		modhash, err := util.BillyCalcFileHash(mdr.ModFile, m.FS)
		if err != nil {
			mdr.errors = append(mdr.errors, err)
			return fmt.Errorf("While calculating billy mod hash\n%w\n", err)
		}


		dver := sumfile.Version{
			Path: strings.Join([]string{m.Module}, "/"),
			Version: m.Version,
		}
		if mdr.module.SumFile == nil {
			mdr.module.SumFile = &sumfile.Sum{}
		}
		mdr.module.SumFile.Add(dver, dirhash)

		mver := sumfile.Version{
			Path: strings.Join([]string{m.Module}, "/"),
			Version: strings.Join([]string{m.Version, mdr.ModFile}, "/"),
		}
		mdr.module.SumFile.Add(mver, modhash)

		baseDir := path.Join(mdr.ModsDir, m.Module)

		fmt.Printf("Writing %-48s => %s\n", m.ReplaceModule + "@" + m.ReplaceVersion, baseDir)

		// copy definite files always
		files, err := m.FS.ReadDir("/")
		if err != nil {
			return err
		}
		for _, file := range files {
			for _, fn := range definiteVendors {
				// Found one!
				if strings.HasPrefix(strings.ToUpper(file.Name()), fn) {
					// TODO, these functions should just take 2 billy FS
					err = util.BillyWriteFileToOS(baseDir, "/"+file.Name(), m.FS)
					if err != nil {
						return err
					}
				}

			}
		}

		if len(mdr.VendorIncludeGlobs) > 0 || len(mdr.VendorExcludeGlobs) > 0 {
			// Just copy everything
			// TODO, these functions should just take 2 billy FS
			err = util.BillyGlobWriteDirToOS(baseDir, "/", m.FS, mdr.VendorIncludeGlobs, mdr.VendorExcludeGlobs)
			if err != nil {
				return err
			}

		} else {
			// Just copy everything
			// TODO, these functions should just take 2 billy FS
			err = util.BillyWriteDirToOS(baseDir, "/", m.FS)
			if err != nil {
				return err
			}

		}

	}

	// Write sumfile
	out, err := mdr.module.SumFile.Write()
	if err != nil {
		return err
	}

	ioutil.WriteFile(mdr.SumFile, []byte(out), 0644)

	return nil
}
