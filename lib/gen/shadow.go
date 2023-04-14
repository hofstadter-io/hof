package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
)

func (G *Generator) CleanupShadow(outputBase, shadowBase string, verbosity int) (errs []error) {
	// calc dirs per generator
	outputDir := filepath.Join(outputBase, G.OutputPath())
	shadowDir := filepath.Join(shadowBase, G.ShadowPath())

	// Cleanup File & Shadow
	errsC := G.CleanupRemainingShadow(outputDir, shadowDir, verbosity)
	errs = append(errs, errsC...)

	// process the subgenerators
	for _, SG := range G.Generators {
		sgerrs := SG.CleanupShadow(outputBase, shadowBase, verbosity)
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

func (G *Generator) LoadShadow(basedir string) (error) {
	shadow := map[string]*File{}

	if G.Verbosity > 1 {
		fmt.Printf("Loading shadow @ %q\n", basedir)
	}

	// make sure the shadow exists
	_, err := os.Lstat(basedir)
	if err != nil {
		// make sure we check err for something actually bad
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" {
			return err
		}
		// file not found, leave politely
		if G.Verbosity > 1 {
			fmt.Println("  shadow not found")
		}
		return nil
	}

	err = filepath.Walk(basedir, func(fpath string, info os.FileInfo, err error) error {
		// Don't need to save directories
		// we should try to clean them up though (todo)
		if info.IsDir() {
			return nil
		}

		// check if first filepath component matches a sub-generator name
		// we need to skip these because the shadow is nested
		// TODO, we could get a conflict if the parent gen writes to a dir with same name as the subgen
		if len(G.Generators) > 0 {
			for _, sg := range G.Generators {
				if G.Verbosity > 2 {
					fmt.Println("checking:", filepath.Join(basedir, sg.Name, "**", "*"), fpath)
				}
				match, err := zglob.Match(filepath.Join(basedir, sg.Name, "**", "*"), fpath)
				if err != nil {
					return err
				}
				if match {
					return nil
				}
			}
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
		if G.Verbosity > 1 {
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
		return err
	}

	G.Shadow = shadow

	return nil
}

// TODO, how to cleanup if the user deletes a generator
