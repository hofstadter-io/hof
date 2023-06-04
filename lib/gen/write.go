package gen

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func (G *Generator) GenerateFiles(outdir string) []error {
	errs := []error{}

	start := time.Now()

	for _, F := range G.OrderedFiles {
		if F.Filepath == "" {
			F.IsSkipped = 1
			continue
		}

		// late bind shadow file to File, because we also late load the shadow dir
		F.ShadowFile = G.Shadow[F.Filepath]

		// TODO, merge G.In with F.In (this should be happening, cannot find)

		// this handles the diff logic
		err := F.Render(outdir, G.UseDiff3, G.NoFormat)
		if err != nil {
			F.IsErr = 1
			F.Errors = append(F.Errors, err)
			errs = append(errs, err)
		}
	}

	elapsed := time.Now().Sub(start).Round(time.Millisecond)
	G.Stats.RenderingTime = elapsed

	return errs
}

func (G *Generator) Write(outputBase, shadowBase string) (errs []error) {
	if G.Disabled {
		return nil
	}

	outputDir := filepath.Join(outputBase, G.OutputPath())
	shadowDir := filepath.Join(shadowBase, G.ShadowPath())

	// TODO, thoughts from thinking about generator monorepos and template/partial/static lookup
	// can we just figure out if the import is the same module by asking CUE?
	// runtime could hold this info, generators declare and add subdir to end (or separate)?
	// do we have to figure out if we can walk cue import tree?
	// and not a prefix of the current module
	// can we extract this base path logic out into a function on the generator?


	writestart := time.Now()

	// Order is important here for implicit overriding of content

	// TODO, separate this into load / write parts
	// treat static like files, can just call render with static=true
	// helpful to load early so we can print file lists of what would be generated, etc...

	// Finally write the generator files
	for _, F := range G.OrderedFiles {
		if G.Verbosity > 1 {
			fmt.Println("Writing:", F.Filepath)
		}

		F.Filepath = filepath.Clean(F.Filepath)
		// Write the actual output
		if F.DoWrite && len(F.Errors) == 0 {
			// todo, lift this out?
			err := F.WriteOutput(outputDir)
			if err != nil {
				errs = append(errs, err)
			}
		}

		// Write the shadow too, or if it doesn't exist
		if G.UseDiff3 {
			if F.DoWrite || (F.IsSame > 0 && F.ShadowFile == nil) {
				err := F.WriteShadow(shadowDir)
				if err != nil {
					errs = append(errs, err)
				}
			}
		}

		// remove from shadows map so we can cleanup what remains
		delete(G.Shadow, F.Filepath)
	}

	// capture timing
	writeend := time.Now()
	G.Stats.WritingTime = writeend.Sub(writestart).Round(time.Millisecond)

	// process the subgenerators
	for _, SG := range G.Generators {
		SG.UseDiff3 = G.UseDiff3
		sgerrs := SG.Write(outputBase, shadowBase)
		errs = append(errs, sgerrs...)
	}

	// capture timing total timing again? (with subgens included)

	// dangit, we do need to account for generators in shadow
	// what if I only want to run one of N generators? (which all output to same dir)

	return errs
}


func (F *File) WriteOutput(basedir string) error {
	// add newline to user output
	F.FinalContent = append(F.FinalContent, '\n')

	// print to stdout
	if F.Filepath == "-" || strings.HasPrefix(F.Filepath, "hof-stdout-") {
		fmt.Print(string(F.FinalContent))
		return nil
	}

	// write to file
	err := F.write(basedir, F.FinalContent)
	if err != nil {
		return err
	}

	F.IsWritten = 1

	return nil
}

func (F *File) WriteShadow(basedir string) error {
	return F.write(basedir, F.RenderContent)
}

func (F *File) write(basedir string, content []byte) error {

	fn := filepath.Join(basedir, F.Filepath)
	dir := filepath.Dir(fn)

	// fix for absolute paths
	if strings.HasPrefix(F.Filepath, "/") {
		_, fn = filepath.Split(F.Filepath)
		fn = filepath.Join(basedir, fn)
		dir = filepath.Dir(fn)
	}

	err := yagu.Mkdir(dir)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// fmt.Println("write", dir, fn, string(content))
	err = ioutil.WriteFile(fn, content, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
