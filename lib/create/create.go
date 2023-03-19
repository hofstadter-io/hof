package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/datautils/io"
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

func Create(module string, extra []string, rootflags flags.RootPflagpole, cmdflags flags.CreateFlagpole) (err error) {
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
	if ver == "" {
		ver = "latest"
	}

	ref := ver
	// ensure we have the most up-to-date code
	if looksLikeRepo(url) {
		_, err = cache.FetchRepoSource(url, ver)
		if err != nil {
			return err
		}
		ref, err = cache.UpgradePseudoVersion(url, ver)
		if err != nil {
			return err
		}
	}


	fmt.Println("setting up...")
	tmpdir, subdir, err = setupTmpdir(url, ref)
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

	// calculate workdir
	outdir := cmdflags.Outdir
	if outdir == "" {
		outdir = cwd
	} else if !filepath.IsAbs(outdir) {
		outdir = filepath.Join(cwd, outdir)
	}
	// fmt.Println("  outdir: ", outdir)

	/*
	// from were we run to root
	rel, err := filepath.Rel(workdir, "/")
	if err != nil {
		// fmt.Println("got here", err, wd, workdir, cwd, cmdflags.Outdir)
		return err
	}
	// fmt.Println("  rel: ", rel)

	// we want a relative path input, the runtime/generator will combine & clean this up
	outdir = filepath.Join(rel, cwd, cmdflags.Outdir)
	fmt.Println("  outdir: ", outdir)
	*/
	genflags.Outdir = outdir

	// create our runtime now, maybe we want a new func for this
	//   since we want to ignore any current CUE module context
	//   everything is put into a temp dir and rendered to CWD
	R, err := gen.NewRuntime(nil, rootflags, genflags)
	if err != nil {
		return err
	}
	R.OriginalWkdir = cwd

	if R.Verbosity > 0 {
		fmt.Println("CueDirs:", R.CueModuleRoot, R.WorkingDir)
	}

	// fmt.Println("pre-run-creator")
	err = runCreator(R, extra, cmdflags.Input)
	// fmt.Println("post-run-creator")
	return err
}

func setupTmpdir(url, ver string) (tmpdir, subdir string, err error) {
	var FS billy.Filesystem

	// put this here for testing and systems where the tmpdir done
	// assuming this is essentially a no-op when it already exists
	err = os.MkdirAll(os.TempDir(), 0755)
	if err != nil {
		return "", "", err
	}

	tmpdir, err = os.MkdirTemp("", "hof-create-")
	if err != nil {
		// fmt.Println("got here", os.TempDir(), err)
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
	out, err := yagu.Shell("hof mod tidy", tmpdir)
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

func runCreator(R *gen.Runtime, extra, inputs []string) (err error) {

	// minimally load and extract generators
	err = R.LoadCue()
	if err != nil {
		return err
	}

	// extract generators
	err = R.ExtractGenerators()
	if err != nil {
		return err
	}
	if len(R.Generators) == 0 {
		return fmt.Errorf("no generators found, please make sure there is a creator at the root of the repository")
	}
	if len(R.Generators) > 1 {
		fmt.Println("Warning, you are running more than one generator. Use --list and -G if this was not your intention.")
	}

	var inputMap map[string]any
	if len(inputs) > 0 {
		// load inputs
		inputMap, err = loadCreateInputs(R, inputs)
		if err != nil {
			return err
		}
		if R.Verbosity > 0 {
			fmt.Println("Create flag-input:", inputMap)
		}
	}

	// handle create input / prompt
	for _, G := range R.Generators {
		err = handleGeneratorCreate(G, extra, inputMap)
		if err != nil {
			return err
		}
	}

	// fmt.Printf("map: %#+v\n", inputMap)

	// full load generators
	// so author can use create inputs in other fields
	errs := R.LoadGenerators()
	if len(errs) > 0 {
		fmt.Println("hello")
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

	// we wait until the very end of all generators to print after messages
	for _, G := range R.Generators {
		after := G.CueValue.LookupPath(cue.ParsePath("Create.Message.After"))
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

func loadCreateInputs(R *gen.Runtime, inputFlags []string) (input map[string]any, err error) {
	if len(inputFlags) == 0 {
		return nil, nil
	}

	input = make(map[string]any)

	for _, inFlag := range inputFlags {
		// starts with @, load file
		// only one time supported right now
		if strings.HasPrefix(inFlag, "@") {
			// this still might not be good enough
			// we may need to remember the original working directory on the runtime
			fn := filepath.Join(R.OriginalWkdir, inFlag[1:])
			// fmt.Println("file flat:", inFlag, fn)

			data, err := loadInputsFromFile(R, fn)
			if err != nil {
				return nil, err
			}

			for k,v := range data {
				input[k] = v
			}

			continue
		}

		// otherwise split by =, path=value
		parts := strings.Split(inFlag, "=")
		if len(parts) != 2 {
			return input, fmt.Errorf("input flag must have 'path=value' format")
		}
		// todo, how to deal with types besides strings (list, int, bool)
		path, value := parts[0], parts[1]
		input[path] = value
		// we'd really prefer to support this, but getting errors from CUE about different runtimes
		// input = input.FillPath(cue.ParsePath(path), value)
	}

	// fmt.Printf("pre-input: %#v\n", input)

	return input, nil
}


func handleGeneratorCreate(G *gen.Generator, extraArgs []string, inputMap map[string]any) (err error) {
	// fill any extra args into generator value
	G.CueValue = G.CueValue.FillPath(cue.ParsePath("Create.Args"), extraArgs)

	genVal := G.CueValue

		// pritn the before message if set, otherwise default
	before := genVal.LookupPath(cue.ParsePath("Create.Message.Before"))
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

	if inputMap != nil {
		// if the user provides a schema for input
		inputVal := genVal.LookupPath(cue.ParsePath("Create.Input"))
		if inputVal.Exists() && inputVal.Err() != nil {
			return inputVal.Err()
		}

		// remake map with types based on schema
		newMap := make(map[string]any)
		for k, v := range inputMap {
			// get the current input val
			ival := inputVal.LookupPath(cue.ParsePath(k))
			if ival.Exists() {
				switch t := v.(type) {
					
					// only handling string inputs
					case string:
						// switch 2
						switch ival.IncompleteKind() {
						// another default copy over
						case cue.StringKind:
							newMap[k] = v

						// interseting part where we convert values
						case cue.BoolKind:
							fmt.Println("boolkind")
							n, err := strconv.ParseBool(t)
							if err != nil {
								return err
							}
							newMap[k] = n
							
						case cue.IntKind:
							n, err := strconv.ParseInt(t, 0, 64)
							if err != nil {
								return err
							}
							newMap[k] = n
							
						case cue.FloatKind:
							n, err := strconv.ParseFloat(t, 64)
							if err != nil {
								return err
							}
							newMap[k] = n
						
						// end interesting inputs
						
						default:
							newMap[k] = v
						}
					default:
						newMap[k] = v
				}
			} else {
				newMap[k] = v
			}
		}

		// fmt.Printf("newMap: %#v\n", newMap)

		G.CueValue = G.CueValue.FillPath(cue.ParsePath("Create.Input"), newMap)
	}

	G.CueValue, err = runPrompt(G.CueValue)
	if err != nil {
		return err
	}

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

func loadInputsFromFile(R *gen.Runtime, fn string) (map[string]any, error) {
	ext := filepath.Ext(fn)[1:]
	data := make(map[string]any)

	if ext == "cue" {
		// read CUE file
		content, err := os.ReadFile(fn)
		if err != nil {
			return nil, err
		}

		// compile CUE
		ctx := R.CueRuntime.CueContext
		v := ctx.CompileBytes(content, cue.Filename(fn))
		if v.Err() != nil {
			return nil, v.Err()
		}

		// decode into Go
		err = v.Decode(&data)
		if err != nil {
			return nil, err
		}

	} else {
		// this will read most datafile types
		var d any
		d = make(map[string]any)
		_, err := io.ReadFile(fn, &d)
		if err != nil {
			return nil, err
		}
		data = d.(map[string]any)
	}
	

	/*
			var data interface{}
			data = make(map[string]any)
			_, err := io.ReadFile(fn, &data)
			if err != nil {
				return input, err
			}
			// fmt.Println("(todo) input: ", fn, data)

			for k,v := range data.(map[string]any) {
				input[k] = v
			}
	*/
	return data, nil
}
