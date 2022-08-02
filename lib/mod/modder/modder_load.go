package modder

import (
	"fmt"

	"github.com/go-git/go-billy/v5/osfs"
)

/* Reads the module files relative to the supplied dir from local FS
- ModFile
- SumFile
- MappingFile
*/

/*
func (mdr *Modder) LoadRootFromFS() error {
	// Load the root module
	err := mdr.LoadMetaFromFS(".")
	if err != nil {
		return err
	}

	return nil
}

func (mdr *Modder) LoadModderFromFS(dir string) error {
	// Load the root module
	err := mdr.LoadMetaFromFS(dir)
	if err != nil {
		return err
	}

	return nil
}

func (mdr *Modder) LoadIndexDepsFromFS(dir string) error {
	// Load the root module
	err := mdr.LoadMetaFromFS(".")
	if err != nil {
		return err
	}

	// Load the root module's deps
	err = mdr.LoadRootDeps()
	if err != nil {
		return err
	}
	// Recurse

	return nil
}
*/

func (mdr *Modder) LoadMetaFromBilly() error {
	// Load the root module
	err := mdr.LoadMetaFromFS(".")
	if err != nil {
		return err
	}

	return nil
}

func (mdr *Modder) LoadMetaFromZip() error {
	// Load the root module
	err := mdr.LoadMetaFromFS(".")
	if err != nil {
		return err
	}

	return nil
}

func (mdr *Modder) LoadMetaFromFS(dir string) error {
	// Shortcut for no load modules, forget the reason for no load...
	if mdr.NoLoad {
		return nil
	}

	// Initialize filesystem
	mdr.FS = osfs.New(dir)

	// Initialzie Module related fields
	mdr.module = &Module{
		FS: mdr.FS,
	}
	mdr.depsMap = map[string]*Module{}

	// Load module files
	var err error
	err = mdr.LoadModFile()
	if err != nil {
		return err
	}

	err = mdr.LoadSumFile()
	if err != nil {
		mdr.errors = append(mdr.errors, err)
		return nil
	}

	err = mdr.LoadMappingsFile()
	if err != nil {
		mdr.errors = append(mdr.errors, err)
		return nil
	}

	return nil
}

// Loads the root modules mod file
func (mdr *Modder) LoadModFile() error {
	fn := mdr.ModFile
	m := mdr.module

	err := m.LoadModFile(fn, false /* Do load replace directives! */)
	if err != nil {
		return err
	}

	// no file
	if m.ModFile == nil {
		return fmt.Errorf("no %q file found, directory not initialized for %s", mdr.ModFile, mdr.Name)
	}

	// empty file checks (based on parsed out values)
	if m.ModFile.Module == nil {
		return fmt.Errorf("your %q file appears to be empty", mdr.ModFile)
	}
	if m.ModFile.Module.Mod.Path == "" {
		return fmt.Errorf("your %q file appears to be missing the module clause", mdr.ModFile)
	}

	m.Module = m.ModFile.Module.Mod.Path

	return nil
}

// Loads the root modules sum file
func (mdr *Modder) LoadSumFile() error {
	fn := mdr.SumFile
	m := mdr.module

	err := m.LoadSumFile(fn)
	if err != nil {
		return err
	}

	return nil
}

// Loads the root modules mapping file
func (mdr *Modder) LoadMappingsFile() error {
	fn := mdr.MappingFile
	m := mdr.module

	err := m.LoadMappingFile(fn)
	if err != nil {
		return err
	}

	return nil
}
