package helpers

import (
	"fmt"
	"os"
	"net/url"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/spf13/pflag"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/runtime"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func LoadRuntime(args []string) (*runtime.Runtime, error) {
	// tui.Log("trace", fmt.Sprintf("Panel.loadRuntime.inputs: %v", args))

	// build eval args & flags from the input args
	var (
		rflags flags.RootPflagpole
		cflags flags.EvalFlagpole
	)
	fset := pflag.NewFlagSet("panel", pflag.ContinueOnError)
	flags.SetupRootPflags(fset, &rflags)
	flags.SetupEvalFlags(fset, &cflags)
	fset.Parse(args)
	args = fset.Args()

	// tui.Log("trace", fmt.Sprintf("Panel.loadRuntime.parsed: %v %v", args, rflags))

	R, err := runtime.New(args, rflags)
	if err != nil {
		tui.Log("error", cuetils.ExpandCueError(err))
		return R, err
	}

	err = R.Load()
	if err != nil {
		tui.Log("error", cuetils.ExpandCueError(err))
		return R, err
	}

	return R, nil
}

func LoadValueFromText(content string) (cue.Value, error) {

	ctx := cuecontext.New()
	v := ctx.CompileString(content, cue.Filename("SourceConfig.Text"))

	return v, nil
}

func LoadValueFromFile(filename string) (cue.Value, error) {

	ctx := cuecontext.New()
	b, err := os.ReadFile(filename)
	if err != nil {
		return cue.Value{}, err
	}
	v := ctx.CompileBytes(b, cue.Filename(filename))

	return v, nil
}

func LoadValueFromHttp(fullurl string) (cue.Value, error) {
	// tui.Log("trace", fmt.Sprintf("Panel.loadHttpValue: %s %s", mode, from))

	// rework any cue/play links
	f := fullurl
	if strings.Contains(fullurl, "cuelang.org/play") {
		u, err := url.Parse(fullurl)
		if err != nil {
			tui.Log("error", err)
			return cue.Value{}, err
		}
		id := u.Query().Get("id")
		f = fmt.Sprintf("https://%s/.netlify/functions/snippets?id=%s", u.Host, id)
	}

	// fetch content
	content, err := yagu.SimpleGet(f)
	if err != nil {
		return cue.Value{}, err
	}
	// add the source
	content = "// from: " + fullurl + "\n\n" + content

	// rebuild, TODO, if scope, use that value and scope.Context() here
	ctx := cuecontext.New()
	v := ctx.CompileString(content, cue.InferBuiltins(true))

	return v, nil
}

func LoadValueFromBash(args []string) (cue.Value, error) {

	wd, err := os.Getwd()
	if err != nil {
		return cue.Value{}, err
	}

	script := strings.Join(args, " ")
	out, err := yagu.Bash(script, wd)
	if err != nil {
		return cue.Value{}, err
	}
	out = "// bash: " + " " + out 

	// compile CUE (json, but all json is CUE, which is why we can add a comment)
	ctx := cuecontext.New()
	v := ctx.CompileString(out, cue.InferBuiltins(true))

	return v, nil
}

