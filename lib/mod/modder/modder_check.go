package modder

import (
	"fmt"
	"path"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/lib/mod/parse/sumfile"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func (mdr *Modder) PartitionSumEntries() ([]string, []string, []string, error) {
	present := []string{}
	missing := []string{}
	local := []string{}

	mod := mdr.module
	sf := mod.SumFile
	if sf == nil {
		return nil, nil, nil, fmt.Errorf("No sum file %q for %s, run 'mvs vendor [%s]' to fix.", mdr.SumFile, mdr.Name, mdr.Name)
	}

	for path, R := range mod.SelfDeps {
		// local replace?
		if strings.HasPrefix(R.NewPath, ".") {
			local = append(local, path)
			continue
		}

		ver := sumfile.Version{
			Path:    path,
			Version: R.NewVersion,
		}

		_, ok := sf.Mods[ver]
		if ok {
			present = append(present, path)
		} else {
			missing = append(missing, path)
		}
	}

	return present, missing, local, nil
}

func (mdr *Modder) CompareSumEntryToVendor(R Replace) error {
	sf := mdr.module.SumFile

	// sumfile ver from Replace
	driver := sumfile.Version{
		Path:    R.NewPath,
		Version: R.NewVersion,
	}
	modver := sumfile.Version{
		Path:    R.NewPath,
		Version: path.Join(driver.Version, mdr.ModFile),
	}

	// sumfile dirhash
	d, ok := sf.Mods[driver]
	if !ok {
		merr := fmt.Errorf("Missing module dirhash in sumfile for '%v' from modfile entry '%v'", driver, R)
		// mdr.errors = append(mdr.errors, merr)
		return merr
	}
	sumDirhash := d[0]

	// sumfile modhash
	m, ok := sf.Mods[modver]
	if !ok {
		merr := fmt.Errorf("Missing module modhash in sumfile for '%v' from modfile entry '%v'", modver, R)
		// mdr.errors = append(mdr.errors, merr)
		return merr
	}
	sumModhash := m[0]

	// fmt.Println("SUMHASH", sumDirhash, sumModhash)

	// load vendor copy, oldpath because that is how it will be imported
	rpath := R.OldPath
	if R.OldPath == "" {
		rpath = R.NewPath
	}
	vpath := path.Join(mdr.ModsDir, rpath)
	// fmt.Println("VPATH", vpath)
	FS := osfs.New(vpath)

	// Calc hashes for vendor from billy
	vdrDirhash, err := yagu.BillyCalcHash(FS)
	if err != nil {
		merr := fmt.Errorf("While calculating vendor dirhash for '%v' from '%v'\n%w\n", driver, R, err)
		// mdr.errors = append(mdr.errors, merr)
		return merr
	}

	vdrModhash, err := yagu.BillyCalcFileHash(mdr.ModFile, FS)
	if err != nil {
		merr := fmt.Errorf("While calculating vendor modhash for '%v' from '%v'\n%w\n", modver, R, err)
		// mdr.errors = append(mdr.errors, merr)
		return merr
	}
	// fmt.Println("VDRHASH", vdrDirhash, vdrModhash)

	mismatch := false
	if sumDirhash != vdrDirhash {
		// merr := fmt.Errorf("Mismatched dir hashes in sumfile for '%v' from modfile entry '%v'", driver, R)
		// mdr.errors = append(mdr.errors, merr)
		mismatch = true
	}
	if sumModhash != vdrModhash {
		// merr := fmt.Errorf("Mismatched modfile hashes in sumfile for '%v' from modfile entry '%v'", modver, R)
		// mdr.errors = append(mdr.errors, merr)
		mismatch = true
	}

	if mismatch {
		return fmt.Errorf("Errors with vendor integrity")
	}

	return nil
}

func (mdr *Modder) CompareLocalReplaceToVendor(R Replace) error {

	// load both into billy
	LFS := osfs.New(R.NewPath)
	VFS := osfs.New(path.Join(mdr.ModsDir, R.OldPath))

	// Calc hashes for replace from billy
	localDirhash, err := yagu.BillyGlobCalcHash(LFS, mdr.VendorIncludeGlobs, mdr.VendorExcludeGlobs)
	if err != nil {
		merr := fmt.Errorf("While calculating local dirhash for '%#+v'\n%w\n", R, err)
		mdr.errors = append(mdr.errors, merr)
		return merr
	}

	localModhash, err := yagu.BillyCalcFileHash(mdr.ModFile, LFS)
	if err != nil {
		merr := fmt.Errorf("While calculating local modhash for '%#+v'\n%w\n", R, err)
		mdr.errors = append(mdr.errors, merr)
		return merr
	}
	// fmt.Println("LCLHASH", localDirhash, localModhash)

	// Calc hashes for vendor from billy
	vdrDirhash, err := yagu.BillyGlobCalcHash(VFS, mdr.VendorIncludeGlobs, mdr.VendorExcludeGlobs)
	if err != nil {
		merr := fmt.Errorf("While calculating vendor dirhash for '%#+v'\n%w\n", R, err)
		mdr.errors = append(mdr.errors, merr)
		return merr
	}

	vdrModhash, err := yagu.BillyCalcFileHash(mdr.ModFile, VFS)
	if err != nil {
		merr := fmt.Errorf("While calculating vendor modhash for '%#+v'\n%w\n", R, err)
		mdr.errors = append(mdr.errors, merr)
		return merr
	}
	// fmt.Println("VDRHASH", vdrDirhash, vdrModhash)

	// Do the check
	mismatch := false
	if localDirhash != vdrDirhash {
		merr := fmt.Errorf("Mismatched dirhash for '%#+v'\n%w\n", R, err)
		mdr.errors = append(mdr.errors, merr)
		mismatch = true
	}
	if localModhash != vdrModhash {
		merr := fmt.Errorf("Mismatched modhash for '%#+v'\n%w\n", R, err)
		mdr.errors = append(mdr.errors, merr)
		mismatch = true
	}

	if mismatch {
		return fmt.Errorf("Errors with vendor integrity")
	}

	return nil
}
