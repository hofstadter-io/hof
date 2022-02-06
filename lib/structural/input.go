package structural

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/json"
	"cuelang.org/go/encoding/yaml"
)

type Input struct {
	Original    string
	Entrypoints []string
	Filename    string
	Filetype    string // yaml, json, cue... toml?
	Expression  string // cue expression to select within document
	Content     []byte
	Value       cue.Value
}

func ReadArg(arg string, ctx *cue.Context, cfg *load.Config) (*Input, error) {
	op, err := ParseInput(arg)
	if err != nil {
		return nil, err
	}
	op, err = LoadInput(op, ctx, cfg)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func ReadGlobs(globs []string, ctx *cue.Context, cfg *load.Config) ([]*Input, error) {
	ins := []*Input{}
	for _, glob := range globs {
		in, err := ReadArg(glob, ctx, cfg)
		if err != nil {
			return nil, err
		}

		ins = append(ins, in)
	}

	return ins, nil
}

// Parses arg into an Input.
// arg can be a value, filename, glob, or - for stdin
// can be <arg>@<expr> to subselect from the root value
// can be <entyrpoint>,<endtrypoint>[@<expr>] to support CUE like args
func ParseInput(arg string) (*Input, error) {
	i := &Input{Original: arg, Filename: arg}

	// does the arg look like a file or a CUE value?
	// this is an overly simple check, but should be sufficient for all formats (CUE, JSON, Yaml)
	if strings.ContainsAny(arg, "{}:") || arg == "_" {
		i.Filename = "expression"
	}

	// look for expression
	if strings.Contains(i.Filename, "@") {
		parts := strings.Split(arg, "@")
		if len(parts) != 2 {
			return i, fmt.Errorf("more than on '@' found for input %q", i.Original)
		}
		i.Filename, i.Expression = parts[0], parts[1]
	}
	// look for entrypoints
	if strings.Contains(i.Filename, ",") {
		i.Entrypoints = strings.Split(i.Filename, ",")
		i.Filename = ""
	}

	return i, nil
}

// Loads a parsed Input and sets the Content and Value
func LoadInput(i *Input, ctx *cue.Context, cfg *load.Config) (*Input, error) {
	if i.Filename == "expression" {
		// probably need to load into overlay
		i.Content = []byte(i.Original)
		i.Value = ctx.CompileString(i.Original)
		return i, i.Value.Err()
	}

	// handle entrypoints
	if i.Entrypoints == nil {
		i.Entrypoints = []string{i.Filename}
	}
	v, err := LoadCueInputs(i.Entrypoints, ctx, cfg)
	if err != nil {
		return i, err
	}
	i.Value = v

	if i.Value.Err() != nil {
		return i, i.Value.Err()
	}

	if i.Expression != "" {
		i.Value = i.Value.LookupPath(cue.ParsePath(i.Expression))
	}

	return i, nil
}

// Loads the entrypoints using the context provided
// returns the value from the load after validating it
func LoadCueInputs(entrypoints []string, ctx *cue.Context, cfg *load.Config) (cue.Value, error) {

	if cfg == nil {
		cfg = &load.Config{
			DataFiles: true,
		}
	}

	bis := load.Instances(entrypoints, cfg)

	bi := bis[0]
	// check for errors on the instance
	// these are typically parsing errors
	if bi.Err != nil {
		return cue.Value{}, bi.Err
	}

	// add back any orphaned files (json/yaml)
	for _, f := range bi.OrphanedFiles {
		d, err := os.ReadFile(f.Filename)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}

		switch f.Encoding {

		case "json":
			A, err := json.Extract(f.Filename, d)
			if err != nil {
				fmt.Println("error: ", err)
				os.Exit(1)
			}

			F := &ast.File{
				Filename: f.Filename,
				Decls:    []ast.Decl{A},
			}
			bi.AddSyntax(F)

		case "yaml":
			F, err := yaml.Extract(f.Filename, d)
			if err != nil {
				fmt.Println("error: ", err)
				os.Exit(1)
			}
			bi.AddSyntax(F)

		default:
			fmt.Println("unknown encoding for", f.Filename, f.Encoding)
		}
	}

	// Use cue.Context to turn build.Instance to cue.Instance
	value := ctx.BuildInstance(bi)
	if value.Err() != nil {
		return cue.Value{}, value.Err()
	}

	// Validate the value
	err := value.Validate(
		cue.ResolveReferences(false),
		cue.Concrete(false),
		cue.Definitions(true),
		cue.Hidden(true),
		cue.Optional(true),
		cue.Attributes(true),
		cue.Docs(false),
	)
	if err != nil {
		return cue.Value{}, err
	}

	return value, nil
}

func LoadGlobs(globs []string) ([]Input, error) {
	// no globs, then stdin
	if len(globs) == 0 {
		globs = []string{"-"}
	}

	// handle special stdin case
	if len(globs) == 1 && globs[0] == "-" {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		i := []Input{Input{Filename: "-", Content: b}}
		return i, nil
	}

	// handle general case
	// we will load into a map to remove duplicates
	// and then extract and sort in a slice
	inputs := make(map[string]Input)
	for _, g := range globs {
		// need to check for expression syntax here

		matches, err := filepath.Glob(g)
		if err != nil {
			return nil, err
		}

		for _, m := range matches {
			// continue on duplicate
			if _, ok := inputs[m]; ok {
				continue
			}

			d, err := os.ReadFile(m)
			if err != nil {
				return nil, err
			}

			// handle input types
			ext := filepath.Ext(m)
			switch ext {
			case ".yml", ".yaml":
				s := fmt.Sprintf(yamlMod, string(d))
				d = []byte(s)
			}

			inputs[m] = Input{Filename: m, Content: d}
		}
	}

	ret := make([]Input, 0)
	for _, i := range inputs {
		ret = append(ret, i)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Filename < ret[j].Filename
	})

	return ret, nil
}

const yamlMod = `
import "encoding/yaml"
#content: #"""
%s
"""#
yaml.Unmarshal(#content)
`
