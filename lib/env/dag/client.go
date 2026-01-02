package dag

import (
	"context"
	"fmt"
	"sync"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type Dag struct {
	mx  sync.RWMutex
	ctx context.Context
	dag *dagger.Client

	// should probably just pass in the flags...
	noCache bool

	// catalog allows us to consolidate references across CUE that might get duplicated
	// as well as each entry holding the value, config, and go types for the entire life-cycle
	cat catalog
	hdl stepHandlerMap
}

func NewClient(ctx context.Context, client *dagger.Client) (d *Dag, err error) {
	d = &Dag{
		ctx: ctx,
		dag: client,
		cat: newCatalog(),
	}

	d.hdl = d.makeStepHandlers()

	return d, nil
}

// probably want to expand this to...
// 1. track more than host stuff (containers, dirs, services, etc...)
// 2. track the config, object, id, usage?
// 3. look up by name or id? (need to fill back the id?) or can we create an interface Key() <- better probably
// 4. enough info to track, walk, schedule, and visualize the dag
// 5. an interface or real type here, instead of any
type catalog map[Keyer]any

func newCatalog() catalog {
	return make(map[Keyer]any)
}

type Keyer interface {
	Key() string
}

// preferences #hof: metadata: [memo|id|name]
func vegMemoKey(e *env.Env) string {
	if e.Hof.Memo != "" {
		return e.Hof.Memo
	}
	if e.Hof.ID != "" {
		return e.Hof.ID
	}
	if e.Hof.Metadata.Name != "" {
		return e.Hof.Metadata.Name
	}
	return ""
}

type StepKind struct {
	Kind  string `json:"$kind"`
	Value cue.Value
}

type Step map[string]any

type stepHandler func(c *dagger.Container, step cue.Value) (*dagger.Container, error)
type stepHandlerMap map[string]stepHandler

func (d *Dag) makeStepHandlers() stepHandlerMap {
	return stepHandlerMap{
		// steps, not cataloged like #things
		// command.cue/go
		// #Cmd
		// #Task

		// container.cue/go
		// #Container
		// #DockerBuild

		// envshh.cue
		"envVar":  d.stepEnvVarHandler,
		"envFile": d.stepEnvFileHandler,
		// #Secret
		"secretVar":  d.stepSecretVarHandler,
		"secretVars": d.stepSecretVarsHandler,

		// exec.cue/go
		"exec": d.stepExecHandler,
		// Script, Sh, Bash, Zsh
		"sync":        d.stepSyncHandler,
		"user":        d.stepUserHandler,
		"workdir":     d.stepWorkdirHandler,
		"entrypoint":  d.stepEntrypointHandler,
		"defaultArgs": d.stepDefaultArgsHandler,
		"defaultTerm": d.stepDefaultTermHandler,
		"terminal":    d.stepTerminalHandler,

		// export.cue/go
		// #ExportDir
		// #ExportFile
		// #ExportImageFile
		// #ExportImage
		// #PublishImage
		// #ExportCuefig
		// #ExportDagger

		// filesystem.cue/go
		// #File
		// #Dir
		"mount": d.stepMountHandler,
		"file":  d.stepFileHandler,
		"dir":   d.stepDirHandler,

		// git.cue/go

		// host.cue/go

		// service.cue/go
		"expose":      d.stepExposeHandler,
		"bindService": d.stepBindServiceHandler,

		// space.cue/go

		// template.cue/go

		// volume.cue/go

		"temp": d.stepTempHandler,
	}
}

type kinder struct {
	Kind string `json:"$kind"`
}

func (d *Dag) Container(val cue.Value, noCache bool) (*dagger.Container, error) {
	d.noCache = noCache

	// it's probably wrong to assume this in general
	var k kinder
	err := val.Decode(&k)
	if err != nil {
		return nil, err
	}

	switch k.Kind {
	case "#container":
		return d.HashContainer(val)
	case "#hostImage":
		return d.HashHostImage(val)
	case "#dockerBuild":
		return d.HashDockerBuild(val)
	default:
		return nil, fmt.Errorf("unsupported build target(%s): %v", k.Kind, val)
	}
}

func (d *Dag) Service(val cue.Value, noCache bool) (*dagger.Service, *hashServiceConfig, error) {
	d.noCache = noCache

	// it's probably wrong to assume this in general
	var k kinder
	err := val.Decode(&k)
	if err != nil {
		return nil, nil, err
	}

	switch k.Kind {
	case "#service":
		s, cfg, err := d.HashService(val)

		return s, cfg, err
	default:
		return nil, nil, fmt.Errorf("unsupported build target(%s): %v", k.Kind, val)
	}
}

func (d *Dag) File(val cue.Value, noCache bool) (*dagger.File, string, error) {
	d.noCache = noCache

	// it's probably wrong to assume this in general
	var k kinder
	err := val.Decode(&k)
	if err != nil {
		return nil, "", err
	}

	switch k.Kind {
	case "#file":
		return d.hashFile(val)
	case "#hostFile":
		file, cfg, err := d.HashHostFile(val)
		if err != nil {
			return nil, "", err
		}
		return file, cfg.Path, nil

	default:
		return nil, "", fmt.Errorf("unsupported build target(%s): %v", k.Kind, val)
	}
}

func (d *Dag) Dir(val cue.Value, noCache bool) (*dagger.Directory, string, error) {
	d.noCache = noCache

	// it's probably wrong to assume this in general
	var k kinder
	err := val.Decode(&k)
	if err != nil {
		return nil, "", err
	}

	switch k.Kind {
	case "#dir":
		return d.hashDir(val)
	case "#hostDir":
		dir, cfg, err := d.HashHostDir(val)
		if err != nil {
			return nil, "", err
		}
		return dir, cfg.Path, nil
	case "#gitRepo":
		repo, rcfg, err := d.hashGitRepo(val)
		if err != nil {
			return nil, "", err
		}
		return repo.Ref(rcfg.Ref).Tree(), "", nil
	default:
		return nil, "", fmt.Errorf("unsupported build target(%s): %v", k.Kind, val)
	}
}
