package envs

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
)

// so we don't have to pass these around everywhere
type Runtime struct {
	Ctx    context.Context
	Client *dagger.Client

	HofVer string
	CueVer string
	GoVer  string

	ContainerRuntime string
	ContainerVersion string

	RunGroup string
	RunCase  string

	Verbose bool
}

func NewRuntime() (*Runtime, error) {
	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return &Runtime{
		Ctx:    ctx,
		Client: client,
	}, nil
}

func (R *Runtime) HofSource() (*dagger.Directory) {

	if R.HofVer == "local" {
		// TODO, improve the dir arg here, right now it requires we run this file from the repo root
		return R.Client.Host().Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"cue.mod/pkg", "docs", "next"},
		})
	} else {
		return R.GitSource("https://github.com/hofstadter-io/hof", R.HofVer)
	}
}
func (R *Runtime) GitSource(url, ref string) (*dagger.Directory) {
		// get repo
		git := R.Client.Git(url, dagger.GitOpts{ KeepGitDir: true })

		// get ref
		var dref *dagger.GitRef
		if strings.HasPrefix(ref, "v") {
			dref = git.Tag(ref)
		} else if strings.HasPrefix(ref, "branch/") {
			dref = git.Branch(strings.TrimPrefix(ref, "branch/"))
		} else {
			dref = git.Commit(ref)
		}

		// return source
		return dref.Tree()
}

func (R *Runtime) BaseImage() (*dagger.Container) {

	golang := R.Client.Container().From("golang:" + R.GoVer)
	c := golang

	// setup mod cache
	modCache := R.Client.CacheVolume("gomod")
	c = c.WithMountedCache("/go/pkg/mod", modCache)

	// setup build cache
	buildCache := R.Client.CacheVolume("go-build")
	c = c.WithMountedCache("/root/.cache/go-build", buildCache)

	// add container tools
	switch R.ContainerRuntime {
	case "docker":
		c = R.AddDockerCLI(c)
	default:
		panic("unsupported runtime: " + R.ContainerRuntime)
	}

	// add packages
	c = c.WithExec([]string{
		"apt-get", "update", "-y",
	})
	c = c.WithExec([]string{
		"apt-get", "install", "-y", 
		"gcc",
		"git",
		"make",
		"python3",
		"tar",
		"tree",
		"wget",
	})

	// get CUE binary
	url := fmt.Sprintf("https://github.com/cue-lang/cue/releases/download/%s/cue_%s_linux_amd64.tar.gz", R.CueVer, R.CueVer)
	tar := fmt.Sprintf("cue_%s_linux_amd64.tar.gz", R.CueVer)
	t := golang.WithWorkdir("/tmp")
	t = t.WithExec([]string{ "wget", url})
	t = t.WithExec([]string{ "ls", "-lh"})
	t = t.WithExec([]string{ "tar", "-xf", tar})
	cue := t.File("/tmp/cue")

	// add CUE binary
	c = c.WithFile("/usr/local/bin/cue", cue)

	// setup workdir
	c = c.WithWorkdir("/work")

	return c
}

func (R *Runtime) FetchGoDeps(c *dagger.Container, source *dagger.Directory) (*dagger.Container) {
	c = c.Pipeline("hof/deps")

	// get deps
	c = c.WithDirectory("/work", source, dagger.ContainerWithDirectoryOpts{
		Include: []string{"go.mod", "go.sums"},
	})
	c = c.WithExec([]string{"go", "mod", "download"})

	// c = c.WithDirectory("/work", source)
	return c
}

func (R *Runtime) BuildHof(c *dagger.Container, source *dagger.Directory) (*dagger.Container) {
	c = c.Pipeline("hof/build")

	// exclude files we don't need so we can avoid cache misses?
	c = c.WithDirectory("/work", source, dagger.ContainerWithDirectoryOpts{
		Exclude: []string{
			"changelogs",
			"ci",
			"docs",
			"hack",
			"images",
			"notes",
			"test", 
		},
	})

	c = c.WithEnvVariable("CGO_ENABLED", "0")

	c = c.WithExec([]string{"go", "build", "./cmd/hof"})
	return c
}

func (R *Runtime) HofImage(builder *dagger.Container) (*dagger.Container) {
	hof := builder.File("hof")

	c := R.BaseImage()
	c = c.WithFile("/usr/local/bin/hof", hof)
	c = c.Pipeline("hof/runtime")
	
	return c
}

