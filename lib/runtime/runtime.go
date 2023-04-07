package runtime

import (
	"os"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/datamodel"
	"github.com/hofstadter-io/hof/lib/hof"
)

// This is the hof Runtime that backs most commands
type Runtime struct {
	// original flags used to load the CUE
	Flags flags.RootPflagpole

	// TODO, can we embed all the command flags here?
	// depending on which command was run
	// is there a dependency injection method (like google/wire)
	// or something like how we dealt with $hof & DHof

	// Other important dirs when loading templates (auto set)
	WorkingDir    string
	CueModuleRoot string
	rootToCwd     string  // module root -> working dir (foo/bar)
	cwdToRoot     string  // module root <- working dir (../..)

	// CUE related fields
	Entrypoints    []string
	CueContext     *cue.Context
	CueConfig      *load.Config
	BuildInstances []*build.Instance
	CueErrors      []error
	FieldOpts      []cue.Option

	// when CUE entrypoints have @placement
	origEntrypoints []string

	// when a user supplies an data.json@path.to.field
	dataMappings    map[string]string

	// The CUE value after all loading
	Value    cue.Value

	// we need to rethink how we organize the code
	// in each of these packages so we can separate
	// the commands from the types and core logic
	Nodes      []*hof.Node[any]
	Datamodels []*datamodel.Datamodel
	// Generators map[string]*gen.Generator
	// Workflows  map[string]*flow.Flow
}

func New(entrypoints []string, rflags flags.RootPflagpole) (*Runtime, error) {
	cfg := &load.Config{
		ModuleRoot: "",
		Module:     "",
		Package:    "",
		Dir:        "",
		Tags:       rflags.Tags,
		TagVars:    load.DefaultTagVars(),
		Tests:      false,
		Tools:      false,
		DataFiles:  false,
		Overlay:    make(map[string]load.Source),
	}

	// package?
	if rflags.Package != "" {
		cfg.Package = rflags.Package
	}

	// inject env?
	if rflags.InjectEnv {
		for _, e := range os.Environ() {
			parts := strings.Split(e, "=")
			k,v := parts[0], parts[1]
			cfg.TagVars[k] = load.TagVar{
				Func: func() (ast.Expr, error) {
					return ast.NewString(v), nil
				},
			}
		}
	}

	r := &Runtime{
		Flags: rflags,
		Entrypoints: entrypoints,
		CueConfig:   cfg,
		dataMappings: make(map[string]string),
	}
	return r, nil
}
