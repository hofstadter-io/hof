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

	// catalog allows us to consolidate references across CUE that might get duplicated
	// as well as each entry holding the value, config, and go types for the entire life-cycle
	cat *catalog
	hdl stepHandlerMap

	watched []*hashCacheIndex
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
type catalog struct {
	sync.Map
}

func newCatalog() *catalog {
	return &catalog{}
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
		"envVars":    d.stepEnvVarsHandler,
		"envFile":    d.stepEnvFileHandler,
		"envAll":     d.stepEnvAllHandler,
		"secretVars": d.stepSecretVarsHandler,
		"secretFile": d.stepSecretFileHandler,

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

		// diff.cue/go
		"changes":   d.stepChangesHandler,
		"patch":     d.stepPatchHandler,
		"patchFile": d.stepPatchFileHandler,

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
		"mount":  d.stepMountHandler,
		"file":   d.stepFileHandler,
		"dir":    d.stepDirHandler,
		"rootfs": d.stepRootFSHandler,

		// git.cue/go

		// host.cue/go
		"unixSocket": d.stepUnixSocketHandler,

		// service.cue/go
		"expose":      d.stepExposeHandler,
		"bindService": d.stepBindServiceHandler,

		// space.cue/go

		// template.cue/go

		// volume.cue/go

		"temp": d.stepTempHandler,

		"withoutDefaultArgs":    d.stepWithoutDefaultArgsHandler,
		"withoutDirectory":      d.stepWithoutDirectoryHandler,
		"withoutEntrypoint":     d.stepWithoutEntrypointHandler,
		"withoutEnvVariable":    d.stepWithoutEnvVariableHandler,
		"withoutExposedPort":    d.stepWithoutExposedPortHandler,
		"withoutFile":           d.stepWithoutFileHandler,
		"withoutFiles":          d.stepWithoutFilesHandler,
		"withoutLabel":          d.stepWithoutLabelHandler,
		"withoutMount":          d.stepWithoutMountHandler,
		"withoutRegistryAuth":   d.stepWithoutRegistryAuthHandler,
		"withoutSecretVariable": d.stepWithoutSecretVariableHandler,
		"withoutUnixSocket":     d.stepWithoutUnixSocketHandler,
		"withoutUser":           d.stepWithoutUserHandler,
		"withoutWorkdir":        d.stepWithoutWorkdirHandler,
	}
}

type kinder struct {
	Kind string `json:"$kind"`
}

func (d *Dag) Container(val cue.Value, noCache bool) (*dagger.Container, error) {
	var err error
	val, err = d.ResolveShouldi(val)
	if err != nil {
		return nil, err
	}
	if !val.Exists() {
		return nil, nil
	}

	// it's probably wrong to assume this in general
	var k kinder
	d.mx.RLock()
	err = val.Decode(&k)
	d.mx.RUnlock()
	if err != nil {
		return nil, err
	}

	switch k.Kind {
	case "#container":
		return d.HashContainer(val, noCache)
	case "#hostImage":
		return d.HashHostImage(val, noCache)
	case "#dockerBuild":
		return d.HashDockerBuild(val, noCache)
	case "#rootfs":
		dir, err := d.hashRootFS(val, noCache)
		if err != nil {
			return nil, err
		}
		return d.dag.Container().WithRootfs(dir), nil
	default:
		return nil, fmt.Errorf("unsupported build target(%s): %v", k.Kind, val)
	}
}

func (d *Dag) Service(val cue.Value, noCache bool) (*dagger.Service, *hashServiceConfig, error) {
	var err error
	val, err = d.ResolveShouldi(val)
	if err != nil {
		return nil, nil, err
	}
	if !val.Exists() {
		return nil, nil, nil
	}

	// it's probably wrong to assume this in general
	var k kinder
	d.mx.RLock()
	err = val.Decode(&k)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, err
	}

	switch k.Kind {
	case "#service":
		s, cfg, err := d.HashService(val, noCache)

		return s, cfg, err
	default:
		return nil, nil, fmt.Errorf("unsupported build target(%s): %v", k.Kind, val)
	}
}

func (d *Dag) File(val cue.Value, noCache bool) (*dagger.File, string, error) {
	var err error
	val, err = d.ResolveShouldi(val)
	if err != nil {
		return nil, "", err
	}
	if !val.Exists() {
		return nil, "", nil
	}

	// it's probably wrong to assume this in general
	var k kinder
	d.mx.RLock()
	err = val.Decode(&k)
	d.mx.RUnlock()
	if err != nil {
		return nil, "", err
	}

	switch k.Kind {
	case "#file":
		return d.hashFile(val, noCache)
	case "#hostFile":
		file, cfg, err := d.HashHostFile(val, noCache)
		if err != nil {
			return nil, "", err
		}
		return file, cfg.Path, nil
	case "#cuefigSBOM":
		return d.HashCuefigSBOM(val, noCache)

	case "#changes":
		chg, err := d.HashChanges(val, noCache)
		file := chg.AsPatch()
		return file, "", err
	case "#patchFile":
		file, err := d.HashPatchFile(val, noCache)
		return file, "", err

	default:
		return nil, "", fmt.Errorf("unsupported build target(%s): %v", k.Kind, val)
	}
}

func (d *Dag) Dir(val cue.Value, noCache bool) (*dagger.Directory, string, error) {
	var err error
	val, err = d.ResolveShouldi(val)
	if err != nil {
		return nil, "", err
	}
	if !val.Exists() {
		return nil, "", nil
	}

	// it's probably wrong to assume this in general
	var k kinder
	d.mx.RLock()
	err = val.Decode(&k)
	d.mx.RUnlock()
	if err != nil {
		return nil, "", err
	}

	switch k.Kind {
	case "#dir":
		return d.hashDir(val, noCache)
	case "#hostDir":
		dir, cfg, err := d.HashHostDir(val, noCache)
		if err != nil {
			return nil, "", err
		}
		return dir, cfg.Path, nil
	case "#rootfs":
		dir, err := d.hashRootFS(val, noCache)
		return dir, "", err
	case "#gitRepo":
		repo, rcfg, err := d.hashGitRepo(val, noCache)
		if err != nil {
			return nil, "", err
		}
		if rcfg != nil && rcfg.Ref != "" {
			return repo.Ref(rcfg.Ref).Tree(), "", nil
		} else {
			return repo.Head().Tree(), "", nil
		}
	default:
		return nil, "", fmt.Errorf("unsupported build target(%s): %v", k.Kind, val)
	}
}
