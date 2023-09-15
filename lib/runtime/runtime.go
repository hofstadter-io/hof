package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/flow/flow"
	"github.com/hofstadter-io/hof/lib/chat"
	"github.com/hofstadter-io/hof/lib/datamodel"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/hof"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

// This is the hof Runtime that backs most commands
type Runtime struct {
	sync.Mutex

	// original flags used to load the CUE
	Flags flags.RootPflagpole

	// TODO, can we embed all the command flags here?
	// depending on which command was run
	// is there a dependency injection method (like google/wire)
	// or something like how we dealt with $hof & DHof

	// Other important dirs when loading templates (auto set)
	WorkingDir    string
	CueModuleRoot string
	RootToCwd     string  // module root -> working dir (foo/bar)
	CwdToRoot     string  // module root <- working dir (../..)
	// OutputDir     string  // where gen wants to write (tbd, other commands too)
	OriginalWkdir string  // when we need to cd and then output back to this directory (create related, but could expand)

	// CUE related fields
	Entrypoints    []string
	CueContext     *cue.Context
	CueConfig      *load.Config
	BuildInstances []*build.Instance
	FieldOpts      []cue.Option

	// this is a bit hacky, but we use this with vet to validate data (and probably st as well)
	DontPlaceOrphanedFiles bool

	// when CUE entrypoints have @placement
	origEntrypoints []string

	// when a user supplies an data.json@path.to.field
	dataMappings    map[string]string

	// internal bookkeeping
	loadedFiles []string

	// The CUE value after all loading
	Value    cue.Value

	// we need to rethink how we organize the code
	// in each of these packages so we can separate
	// the commands from the types and core logic
	Nodes      []*hof.Node[any]
	Chats      []*chat.Chat
	Datamodels []*datamodel.Datamodel
	Generators []*gen.Generator
	Workflows  []*flow.Flow

	Stats RuntimeStats
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
		origEntrypoints: entrypoints,
		CueConfig:   cfg,
		dataMappings: make(map[string]string),
		Stats: make(RuntimeStats),
	}

	// calc cue dirs
	var err error
	r.CueModuleRoot, err = cuetils.FindModuleAbsPath("")
	if err != nil {
		return r, err
	}
	// TODO: we could make this configurable
	r.WorkingDir, _ = os.Getwd()
	if r.CueModuleRoot != "" {
		r.CwdToRoot, err = filepath.Rel(r.WorkingDir, r.CueModuleRoot)
		if err != nil {
			return r, err
		}
		r.RootToCwd, err = filepath.Rel(r.CueModuleRoot, r.WorkingDir)
		if err != nil {
			return r, err
		}
	}

	return r, nil
}

// OutputDir returns the absolute path to output dir for this runtime.
// It accounts for module root and relative directories.
func (R *Runtime) OutputDir(dir string) string {
	if strings.HasPrefix(dir, "/") {
		return dir
	}
	return filepath.Join(R.CueModuleRoot, R.RootToCwd, dir)
}

func (R *Runtime) GetLoadedFiles() []string {
	var files []string
	bi := R.BuildInstances[0]

	// these two should cover us, though we might need to process imports?
	for _, f := range bi.BuildFiles {
		files = append(files, f.Filename)
	}
	for _, f := range bi.OrphanedFiles {
		files = append(files, f.Filename)
	}

	return files
}
