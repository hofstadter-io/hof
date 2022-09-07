package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
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

func Create(module string, rootflags flags.RootPflagpole, cmdflags flags.CreateFlagpole) (err error) {
	var tmpdir, subdir string

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	parts := strings.Split(module, "@")
	url, ver := parts[0], ""
	if len(parts) == 2 {
		ver = parts[1]
	}
	// fmt.Printf("looking for %s @ %s\n", url, ver)

	fmt.Println("setting up...")
	tmpdir, subdir, err = setupTmpdir(url, ver)
	if err != nil {
		return fmt.Errorf("while setting up tmpdir: %w", err)
	}

	workdir := filepath.Join(tmpdir, subdir)
	// fmt.Printf("dirs: %q %q %q\n", tmpdir, subdir, workdir)

	//
	// chdir to tmpdir....
	//

	err = os.Chdir(workdir)
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)

	// defer jumping back
	defer func() {
		// fmt.Println("cleaning up", tmpdir)
		os.RemoveAll(tmpdir)
	}()

	fmt.Println("init'n generator")
	genflags := flags.GenFlags
	genflags.Generator = cmdflags.Generator
	genflags.Outdir = cmdflags.Outdir

	// from were we run to root
	rel, err := filepath.Rel(workdir, "/")
	if err != nil {
		return err
	}
	// fmt.Println("  rel: ", rel)

	// we want a relative path input, the runtime/generator will combine & clean this up
	outdir := filepath.Join(rel, cwd, cmdflags.Outdir)
	// fmt.Println("  outdir: ", outdir)
	genflags.Outdir = outdir

	// create our runtime now, maybe we want a new func for this
	//   since we want to ignore any current CUE module context
	//   everything is put into a temp dir and rendered to CWD
	R, err := gen.NewRuntime(nil, rootflags, genflags)
	if err != nil {
		return err
	}

	if R.Verbosity > 0 {
		fmt.Println("CueDirs:", R.CueModuleRoot, R.WorkingDir)
	}

	// fmt.Println("pre-run-creator")
	err = runCreator(R, cmdflags.Input)
	// fmt.Println("post-run-creator")
	return err
}

func setupTmpdir(url, ver string) (tmpdir, subdir string, err error) {
	var FS billy.Filesystem

	tmpdir, err = os.MkdirTemp("", "hof-create-")
	if err != nil {
		return tmpdir, "", err
	}
	err = os.MkdirAll(tmpdir, 0755)
	if err != nil {
		return tmpdir, "", err
	}

	// remote, or local
	if looksLikeRepo(url) {
		fmt.Println("loading remote creator")
		// todo, this should handle walking up and cloning
		// ? we need to determine subdir during cache walkback and return?
		FS, err = cache.Load(url, ver)
		if err != nil {
			return tmpdir, "", err
		}

	} else {
		fmt.Println("loading local creator")
		// fmt.Println("local creator")
		// todo, walk up to mod / git root
		// ? how to deal with subdirs
		// check for directory
		info, err := os.Lstat(url)
		if err != nil {
			return tmpdir, "", err
		}

		if !info.IsDir() {
			return tmpdir, "", fmt.Errorf("%s is not a directory", url)
		}

		// find cue module root from input url
		// fmt.Println("url:", url)
		modroot, err := cuetils.FindModuleAbsPath(url)
		if err != nil {
			return tmpdir, "", err
		}

		// fmt.Println("modroot:", modroot)

		// abs path of input for next calc
		abs, err := filepath.Abs(url)
		if err != nil {
			return tmpdir, "", err
		}
		// fmt.Println("abs:", abs)

		// find subdir, after both are absolute
		subdir, err = filepath.Rel(modroot, abs)
		if err != nil {
			return tmpdir, subdir, err
		}

		// fmt.Println("subdir:", subdir)

		// load into FS
		// fmt.Println("starting to read:", modroot)
		FS = osfs.New(modroot)
		// fmt.Println("done reading")

	}

	// fmt.Println("writing", tmpdir)

	err = yagu.BillyWriteDirToOS(tmpdir, "/", FS)
	if err != nil {
		return tmpdir, subdir, err
	}

	// run 'hof mod vendor cue' in tmpdir
	fmt.Println("fetching creator dependencies")
	out, err := yagu.Bash("hof mod vendor cue", tmpdir)
	// fmt.Println("done fetching dependencies\n", out)
	if err != nil {
		fmt.Println(out)
		return tmpdir, subdir, fmt.Errorf("while fetching creator deps %w", err)
	}

	// DEV DEBUG informational only
	/*
	infos, err := os.ReadDir(tmpdir)
	for _, info := range infos {
		fmt.Println(info.Name())
	}
	*/
	
	return tmpdir, subdir, err
}

func runCreator(R *gen.Runtime, inputs []string) (err error) {

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
