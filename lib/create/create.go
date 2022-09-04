package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/repos/cache"
	"github.com/hofstadter-io/hof/lib/yagu"
)

/*  Steps for create

- local or remote?
- copy to tempdir
- run hof mod vendor cue
- run create process
- cleanup

need to get dirs aligned with various tools (filepaths... smh)
will also need to walk up dirs to find git/mod roots

common create prompts

- name
- git repo (setup/confirm)
- cue repo (setup/confirm)
- outdir (depends on others without flag)
  - if in repo, probably write to current dir
  - if not in repo, mkdir and write there
  - what does npm/x create-* do?
  - noting that we may deviate because people may use create to add to their existing projects
*/

func Create(args []string, rootflags flags.RootPflagpole, cmdflags flags.CreateFlagpole) (err error) {

	// tempdir if it gets filled
	tmpdir := ""
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}


	// fmt.Println("Create:", args, cmdflags)
	// mdr := mod.LangModderMap["cue"]
	// mdr.SetWorkdir()

	// TODO, is this local or remote?
	// or is this cue entrypoints or does it look like a remote
	if len(args) == 1 {
		if looksLikeRepo(args[0]) {
			parts := strings.Split(args[0], "@")
			url, ver := parts[0], ""
			if len(parts) == 2 {
				ver = parts[1]
			}
			fmt.Printf("looking for %s @ %s\n", url, ver)

			tmpdir, err = setupTmpdir(url, ver)
			if err != nil {
				return err
			}

			fmt.Println(tmpdir)
			args = []string{}

			//
			// chdir to tmpdir....
			//

			err = os.Chdir(tmpdir)
			if err != nil {
				return err
			}
			defer os.Chdir(cwd)

			// defer jumping back
			defer func() {
				fmt.Println("cleaning up", tmpdir)
				// os.RemoveAll(tmpdir)
			}()
		} else {

			// if current module, just continue

			// if different module

		}
	}

	if tmpdir != "" {
		outdir := filepath.Join("../../", cwd, cmdflags.Outdir)
		fmt.Println("tmp: ", outdir)
	}

	genflags := flags.GenFlags
	genflags.Generator = cmdflags.Generator
	genflags.Diff3 = false

	R, err := gen.NewRuntime(args, rootflags, genflags)
	if err != nil {
		return err
	}


	if R.Verbosity > 0 {
		fmt.Println("CueDirs:", R.CueModuleRoot, R.WorkingDir)
	}

	fmt.Println("pre-run-creator")
	err = runCreator(R)
	fmt.Println("post-run-creator")
	return err
}

func setupTmpdir(url, ver string) (dir string, err error) {
	tmpdir, err := os.MkdirTemp("", "hof-")
	if err != nil {
		return tmpdir, err
	}

	FS, err := cache.Load(url, ver)
	if err != nil {
		return tmpdir, err
	}


	err = yagu.BillyWriteDirToOS(tmpdir, "/", FS)
	if err != nil {
		return tmpdir, err
	}

	// run 'hof mod vendor cue' in tmpdir
	out, err := yagu.Bash("hof mod vendor cue", tmpdir)
	fmt.Println(out)
	if err != nil {
		return tmpdir, err
	}

	infos, err := os.ReadDir(tmpdir)
	for _, info := range infos {
		fmt.Println(info.Name())
	}
	// fmt.Println("Remote repositories are not supported yet")
	
	return tmpdir, err
}

func runCreator(R *gen.Runtime) (err error) {

	// minimally load and extract generators
	err = R.LoadCue()
	if err != nil {
		return err
	}

	err = R.ExtractGenerators()
	if err != nil {
		return err
	}

	if len(R.Generators) > 1 {
		fmt.Println("Warning, you are running more than one generator. Use --list and -G if this was not your intention.")
	}

	// handle create input / prompt
	for _, G := range R.Generators {
		err = handleGeneratorCreate(G)
		if err != nil {
			return err
		}
	}

	// full load generators
	// so author can use create inputs in other fields
	errs := R.LoadGenerators()
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		return fmt.Errorf("While loading generators")
	}

	// run generators
	errs = R.RunGenerators()
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		return fmt.Errorf("While generating")
	}

	// write output
	errs = R.WriteOutput()
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		return fmt.Errorf("While writing")
	}

	for _, G := range R.Generators {
		after := G.CueValue.LookupPath(cue.ParsePath("CreateMessage.After"))
		if after.Err() != nil {
			fmt.Println("error:", after.Err())
			return after.Err()
		}

		if !after.IsConcrete() || !after.Exists() {
			if R.Verbosity > 0 {
				fmt.Println("done creating")
			}
		} else {
			s, err := after.String()
			if err != nil {
				return err
			}
			fmt.Println(s)
		}
	}

	return nil
}

func handleGeneratorCreate(G *gen.Generator) error {
	before := G.CueValue.LookupPath(cue.ParsePath("CreateMessage.Before"))
	if before.Err() != nil {
		fmt.Println("error:", before.Err())
		return before.Err()
	}

	if !before.IsConcrete() || !before.Exists() {
		fmt.Printf("Creating from %q\n", G.Name)
	} else {
		s, err := before.String()
		if err != nil {
			return err
		}
		fmt.Println(s)
	}

	val := G.CueValue.LookupPath(cue.ParsePath("CreateInput"))
	if val.Err() != nil {
		fmt.Println("error:", val.Err())
		return val.Err()
	}

	if !val.IsConcrete() {
		return fmt.Errorf("Generator is missing CreateInput")
	}

	prompt := G.CueValue.LookupPath(cue.ParsePath("CreatePrompt"))
	if prompt.Err() != nil {
		fmt.Println("error:", prompt.Err())
		return prompt.Err()
	}

	if !prompt.IsConcrete() || !prompt.Exists() {
		return fmt.Errorf("Generator is missing CreatePrompt")
	}

	// fmt.Printf("%s: %v\n", G.Name, val)
	// fmt.Println(prompt)

	ans := map[string]any{}
	// TODO deal with --input flags

	// process create prompts
	// Loop through all top level fields
	iter, err := prompt.List()
	if err != nil {
		return err
	}
	for iter.Next() {
		value := iter.Value()
		Q := map[string]any{}
		err := value.Decode(&Q)
		if err != nil {
			return err
		}

		// fmt.Println("q:", Q)
		// todo, extract Name
		A, err := handleQuestion(Q)
		if err != nil {
			return err
		}

		// do we want to return a bool from handleQuestion
		// to be more explicit about this check?
		if A != nil {
			ans[Q["Name"].(string)] = A
		}
	}

	// fill CreateInput from --inputs and prompt
	G.CueValue = G.CueValue.FillPath(cue.ParsePath("CreateInput"), ans)

	// fmt.Println("Final:", G.CueValue)
	// return fmt.Errorf("intentional error")
	return nil
}

func looksLikeRepo(str string) bool {
	// todo, check if file exists here when returning false
	// (and assuming a single arg cue entrypoint
	// basically we are trying to eliminae this case and
	// if we can't we assume a repo

	// note, we will have to revisit if we start supporting
	// named create args like npm create-*

	// if contains(@), then should be a repo
	if strings.Contains(str, "@") {
		return true
	}

	// does it have '/' separated parts?
	parts := strings.Split(str, "/")

	// if only 1, no slashes, probably a file or directory
	if len(parts) == 1 {
		return false
	}

	// first part should be a domain

	// is it only '.' or '..', it's relative or a dir
	if parts[0] == "." || parts[0] == ".." || !strings.Contains(parts[0], ".") {
		return false
	}

	// if last part contains a '.' or '*', probably an entrypoint
	last := parts[len(parts)-1]
	if strings.Contains(last, ".") || strings.Contains(last, "*") {
		return false
	}

	// we've checked all we can, assume a repo
	return true
}
