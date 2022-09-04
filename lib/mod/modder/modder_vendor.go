package modder

import (
	"fmt"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/lib/repos/cache"
	"github.com/hofstadter-io/hof/lib/yagu"
)

/* Vendor reads in a module, determines dependencies, and writes out the vendor folder.

Will there be infinite recursion, or maybe just two levels?


This module & deps
	1) load this module
		a) check sum <-> mod files, need to determine if any updates here, for now w/o sum file
	2) process its sum/mod files
	3) for each of this mods deps dependency
	  a) fetch all refs
		b) find minimum
		c) add to depMap, when / if ... ? guessing how right now
		  - replaces are processed first
			- requires are processed second, so only add if not there, we shouldn't have duplicates in owr own mods files
			-
	  d) if added, clond the desired ref to memory

	4) Now loop over depMap to pull in secondary dependencies
	  - probably want to create a "newDeps" map here if we need to support wider recursion
		- basically follow the last block, but load idependently and merge after
		- do we need a separate modder when we process each dep?
		  - probably if we are going to enable each module to optionally specify local behavior
		- so first file we should read is the .mvsconfig, that maps <lang> to whatever

	F) Finally, write out the vendor directory
	  a) check <vendor-dir>/modules.txt and checksums
		b) write out if necessary
*/
func (mdr *Modder) Vendor() error {
	// TODO, run pre vendor commands here

	// Vendor Command Override
	if len(mdr.CommandVendor) > 0 {
		for _, cmd := range mdr.CommandVendor {
			out, err := yagu.Exec(cmd)
			fmt.Println(out)
			if err != nil {
				return err
			}
		}
	} else {
		// Otherwise, MVS venodiring
		err := mdr.VendorMVS()
		if err != nil {
			mdr.PrintErrors()
			return err
		}
	}

	// TODO, run post vendor commands here

	return nil
}

// The entrypoint to the MVS internal vendoring process
func (mdr *Modder) VendorMVS() error {
	var err error

	// Load minimal root module
	err = mdr.LoadMetaFromFS(".")
	if err != nil {
		// fmt.Println(err)
		return err
	}
	for _, R := range mdr.module.SelfDeps {
		err := mdr.VendorDep(R)
		if err != nil {
			mdr.errors = append(mdr.errors, err)
		}
	}

	if err := mdr.CheckForErrors(); err != nil {
		return err
	}

	// Finally, write out anything that needs to be
	err = mdr.WriteVendor()
	if err != nil {
		return err
	}
	return nil
}

func (mdr *Modder) VendorDep(R Replace) error {
	// fmt.Printf("VendorDep %#+v\n", R)

	// Fetch and Load module
	if strings.HasPrefix(R.NewPath, "./") || strings.HasPrefix(R.NewPath, "../") {
		err := mdr.LoadLocalReplace(R)
		if err != nil {
			mdr.errors = append(mdr.errors, err)
			return err
		}
		return nil
	} else {
		err := mdr.LoadRemoteModule(R)
		if err != nil {
			mdr.errors = append(mdr.errors, err)
			return err
		}
		return nil
	}

	return nil
}

func (mdr *Modder) LoadRemoteModule(R Replace) (err error) {
	// fmt.Printf("LoadRemoteReplace %#+v\n", R)

	// If sumfile, check integrity and possibly shortcut
	if mdr.module.SumFile != nil {
		mdr.CompareSumEntryToVendor(R)
	}

	// TODO, check if valid and just add m.deps to processing
	// We can do this in a BFS manner

	m := &Module{
		Module:         R.OldPath,
		Version:        R.OldVersion,
		ReplaceModule:  R.NewPath,
		ReplaceVersion: R.NewVersion,
	}

	if m.Module == "" {
		m.Module = R.NewPath
		m.Version = R.NewVersion
	}

	m.FS, err = cache.Load(R.NewPath, R.NewVersion)
	if err != nil {
		return err
	}

	err = m.LoadMetaFiles(mdr.ModFile, mdr.SumFile, mdr.MappingFile, true /* ignoreReplace directives */)
	if err != nil {
		return err
	}

	// TODO, check if valid and just add m.deps to processing

	// XXX Within this call is what basically makes us BFS
	// It's where we load the dependencies
	err = mdr.MvsMergeDependency(m)
	if err != nil {
		return err
	}

	return nil
}

func (mdr *Modder) LoadLocalReplace(R Replace) error {
	// fmt.Printf("LoadLocalReplace %#+v\n", R)
	var err error

	m := &Module{
		Module:         R.OldPath,
		Version:        R.OldVersion,
		ReplaceModule:  R.NewPath,
		ReplaceVersion: R.NewVersion,
	}

	m.FS = osfs.New(R.NewPath)

	err = m.LoadMetaFiles(mdr.ModFile, mdr.SumFile, mdr.MappingFile, true /* ignoreReplace directives */)
	if err != nil {
		return err
	}

	// TODO, check if valid and just add m.deps to processing

	// XXX Within this call is what basically makes us BFS
	// It's where we load the dependencies
	err = mdr.MvsMergeDependency(m)
	if err != nil {
		return err
	}

	return nil
}
