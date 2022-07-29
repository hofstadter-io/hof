package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const SHADOW_DIR = ".hof/shadow/"

func (R *Runtime) CleanupRemainingShadow() (errs []error) {
	if R.Verbosity > 0 {
		fmt.Println("Cleaning shadow")
	}

	for _, G := range R.Generators {
		gerrs := R.CleanupGeneratorShadow(G)
		errs = append(errs, gerrs...)
	}

	return errs
}

func (R *Runtime) CleanupGeneratorShadow(G *Generator) (errs []error) {
	// calc dirs per generator
	outputDir := filepath.Join(R.OutputDir(), G.OutputPath())
	shadowDir := filepath.Join(R.ShadowDir(), G.ShadowPath())

	// Cleanup File & Shadow
	errsC := G.CleanupRemainingShadow(outputDir, shadowDir, R.Verbosity)
	errs = append(errs, errsC...)

	// process the subgenerators
	for _, SG := range G.Generators {
		sgerrs := R.CleanupGeneratorShadow(SG)
		errs = append(errs, sgerrs...)
	}

	return errs
}

func (G *Generator) CleanupRemainingShadow(outputDir, shadowDir string, verbosity int) (errs []error) {
	// no need if not diff3
	if !G.UseDiff3 {
		return nil
	}

	for f := range G.Shadow {
		genFilename := filepath.Join(outputDir, f)
		shadowFilename := filepath.Join(shadowDir, f)
		if verbosity > 0 {
			fmt.Println("  -", G.Name, f, genFilename, shadowFilename)
		} else {
			fmt.Println("  -", f)
		}

		// TODO (good-first-issue, small)
		// this is a good place to delete dirs
		// (1) get dirname of each
		// (2) if empty, then remove

		err := os.Remove(genFilename)
		if err != nil {
			errs = append(errs, err)
		}

		err = os.Remove(shadowFilename)
		if err != nil {
			errs = append(errs, err)
		}

		G.Stats.NumDeleted += 1
	}

	return errs
}

func LoadShadow(basedir string, verbosity int) (map[string]*File, error) {
	shadow := map[string]*File{}

	if verbosity > 1 {
		fmt.Printf("Loading shadow @ %q\n", basedir)
	}

	// make sure the shadow exists
	_, err := os.Lstat(basedir)
	if err != nil {
		// make sure we check err for something actually bad
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" {
			return nil, err
		}
		// file not found, leave politely
		if verbosity > 1 {
			fmt.Println("  shadow not found")
		}
		return nil, nil
	}

	err = filepath.Walk(basedir, func(fpath string, info os.FileInfo, err error) error {
		// Don't need to save directories
		// we should try to clean them up though (todo)
		if info.IsDir() {
			return nil
		}

		// read contents
		bytes, err := os.ReadFile(fpath)
		if err != nil {
			return err
		}

		// trim ShadowDir so we only have path relative to output dir
		fpath = strings.TrimPrefix(fpath, basedir)
		// should never have slash at beginning
		fpath = strings.TrimPrefix(fpath, "/")

		// debug
		if verbosity > 1 {
			fmt.Println("  adding:", fpath)
		}

		shadow[fpath] = &File{
			FinalContent: bytes,
			Filepath: fpath,
		}

		return nil
	})

	if err != nil {
		err = fmt.Errorf("error walking the shadow dir %q: %w\n", basedir, err)
		return nil, err
	}

	return shadow, nil
}

// TODO, how to cleanup if the user deletes a generator
