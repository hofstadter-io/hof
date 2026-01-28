package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"
	"dagger.io/dagger"
	"gorm.io/gorm"

	"github.com/hofstadter-io/cinful"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/flow/flow"
	"github.com/hofstadter-io/hof/lib/agent"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/datamodel"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/hof"
	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/labstack/echo/v4"
)

// This is the hof Runtime that backs most commands
type Runtime struct {
	sync.Mutex
	Ctx context.Context

	// internal service handlers
	DB  *gorm.DB
	API *echo.Echo
	DAG *dagger.Client

	// original flags used to load the CUE
	Flags flags.RootPflagpole

	// TODO, can we embed all the command flags here?
	// depending on which command was run
	// is there a dependency injection method (like google/wire)
	// or something like how we dealt with $hof & DHof

	// Other important dirs when loading templates (auto set)
	WorkingDir    string
	CueModuleRoot string
	CueExtractDir string // where CUE extract modules to
	RootToCwd     string // module root -> working dir (foo/bar)
	CwdToRoot     string // module root <- working dir (../..)
	// OutputDir     string  // where gen wants to write (tbd, other commands too)
	OriginalWkdir string // when we need to cd and then output back to this directory (create related, but could expand)

	DepMapping map[string]string // map of module paths to module names, used for loading modules

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
	dataMappings map[string]string

	// internal bookkeeping
	loadedFiles []string // cue+data (?)

	// non-data files loaded <cue-path> => <file-path>
	userFiles map[string]string // non-data files loaded by user
	modFiles  map[string]string // non-data files loaded by modules

	// The CUE value after all loading
	Value cue.Value

	// we need to rethink how we organize the code
	// in each of these packages so we can separate
	// the commands from the types and core logic
	Nodes []*hof.Node[any]

	// Chats      []*chat.Chat
	Envs       []*env.Env
	Agentics   []*agent.Agentic
	Datamodels []*datamodel.Datamodel
	Generators []*gen.Generator
	Workflows  []*flow.Flow

	Stats RuntimeStats
}

func New(entrypoints []string, rflags flags.RootPflagpole) (*Runtime, error) {
	cfg := &load.Config{
		ModuleRoot:          "",
		Module:              "",
		Package:             "",
		Dir:                 "",
		Tags:                rflags.Tags,
		TagVars:             load.DefaultTagVars(),
		Tests:               false,
		Tools:               false,
		DataFiles:           false,
		Overlay:             make(map[string]load.Source),
		AcceptLegacyModules: true,
	}

	// package?
	if rflags.Package != "" {
		cfg.Package = rflags.Package
	}

	// some more default tag vars
	extra := make(map[string]string)
	// GIT INFO
	if yagu.InGitRepo() {
		extra["gitRoot"], _ = yagu.GitRepoRoot()
		extra["gitCommit"], _ = yagu.GitCommit()
		extra["gitShortSha"], _ = yagu.GitShortSHA()
		extra["gitBranch"], _ = yagu.GitBranch()
		extra["gitTag"], _ = yagu.GitTag()
		dirty, _, _ := yagu.GitDirty()
		if dirty {
			extra["gitDirty"] = "dirty"
		} else {
			extra["gitDirty"] = ""
		}
	} else {
		extra["gitRoot"] = ""
		extra["gitCommit"] = ""
		extra["gitShortSha"] = ""
		extra["gitBranch"] = ""
		extra["gitTag"] = ""
		extra["gitDirty"] = ""
	}

	// CI
	vendor := cinful.Info()
	if vendor != nil {
		extra["ci"] = fmt.Sprint(vendor)
	} else {
		extra["ci"] = ""
	}
	// add extra to TagVars
	for k, v := range extra {
		cfg.TagVars[k] = load.TagVar{
			Func: func() (ast.Expr, error) {
				return ast.NewString(v), nil
			},
		}
	}

	// inject env?
	if rflags.InjectEnv {
		for _, e := range os.Environ() {
			parts := strings.Split(e, "=")
			k, v := parts[0], parts[1]
			cfg.TagVars[k] = load.TagVar{
				Func: func() (ast.Expr, error) {
					return ast.NewString(v), nil
				},
			}
		}
	}

	r := &Runtime{
		Ctx:             context.Background(),
		Flags:           rflags,
		Entrypoints:     entrypoints,
		origEntrypoints: entrypoints,
		CueConfig:       cfg,
		DepMapping:      make(map[string]string),
		dataMappings:    make(map[string]string),
		userFiles:       make(map[string]string),
		modFiles:        make(map[string]string),
		Stats:           make(RuntimeStats),
	}

	// calc cue dirs
	var err error
	r.CueModuleRoot, err = cuetils.FindModuleAbsPath("")
	if err != nil {
		return r, err
	}

	d, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	r.CueExtractDir = filepath.Join(d, "cue", "mod", "extract")
	// fmt.Println("CueExtractDir:", r.CueExtractDir)

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

const VEG_DAGGER_ENGINE_ENV_VAR = "VEG_EXPERIMENTAL_DAGGER_RUNNER_HOST"
const VEG_DAGGER_HOST = "container://veg-dagger-engine"

func (R *Runtime) DaggerInit() (err error) {
	userVal := os.Getenv(VEG_DAGGER_ENGINE_ENV_VAR)
	if userVal == "" {
		os.Setenv(VEG_DAGGER_ENGINE_ENV_VAR, VEG_DAGGER_HOST)
	}
	R.DAG, err = dagger.Connect(R.Ctx)
	if err != nil {
		return fmt.Errorf("while connecting to dagger: %w", err)
	}
	return nil

}
