package create

import (
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func Create(args []string, rootflags flags.RootPflagpole, cmdflags flags.CreateFlagpole) error {

	// fmt.Println("Create:", args, cmdflags)

	// TODO, is this local or remote?
	// or is this cue entrypoints or does it look like a remote
	if len(args) == 1 && looksLikeRepo(args[0]) {
		parts := strings.Split(args[0], "@")
		url, ver := parts[0], ""
		if len(parts) == 2 {
			ver = parts[1]
		}
		fmt.Printf("looking for %s @ %s\n", url, ver)
		root, tmpdir, err := yagu.FindRemoteRepoRootAndClone(url, ver)
		fmt.Println(root, tmpdir, err)
		if err == nil {
			infos, err := os.ReadDir(tmpdir)
			for _, info := range infos {
				fmt.Println(info.Name())
			}
			err = os.RemoveAll(tmpdir)
			if err != nil {
				return err
			}
		}
		fmt.Println("Remote repositories are not supported yet")
		return nil
	}

	// Do local generators need to be moved to vendor so that
	// template lookups still work???
	// should we be making a temp dir, run from there, and then set output dir to abolute path?

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

	// minimally load and extract generators
	err = R.LoadCue()
	if err != nil {
		return err
	}

	err = R.ExtractGenerators()
	if err != nil {
		return err
	}

	// handle create input / prompt
	for _, G := range R.Generators {
		err = handleGeneratorCreate(G)
		if err != nil {
			return err
		}
	}

	// full load generators
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
