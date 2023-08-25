package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"dagger.io/dagger"
	"github.com/kr/pretty"
	"github.com/spf13/pflag"

	"github.com/hofstadter-io/hof/ci/harmony/envs"
	"github.com/hofstadter-io/hof/ci/harmony/registry"
)

var R *envs.Runtime

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	var err error
	R, err = envs.NewRuntime()
	if err != nil {
		panic(err)
	}

	pflag.StringVar(&R.HofVer, "hof", "local", "set Hof version")
	pflag.StringVar(&R.CueVer, "cue", "v0.6.0", "set CUE version")
	pflag.StringVar(&R.GoVer, "go", "1.21", "set Go version")
	pflag.StringVar(&R.ContainerRuntime, "container-runtime", "docker", "set container runtime")
	pflag.StringVar(&R.ContainerVersion, "container-version", "24", "set container version")
	pflag.StringVar(&R.RunGroup, "group", "", "run tests where the group name has this flag as a prefix")
	pflag.StringVar(&R.RunGroup, "case", "", "run tests where the case name has this flag as a prefix")
	pflag.BoolVar(&R.Verbose, "verbose", false, "increase logging output")
}

const msgfmt = `running harmony
  hof: %s
  cue: %s
  go:  %s
  %s: %s
`

func main() {
	pflag.Parse()
	fmt.Printf(msgfmt, R.HofVer, R.CueVer, R.GoVer, R.ContainerRuntime, R.ContainerVersion)

	// foundation
	source := R.HofSource()
	base := R.BaseImage()

	// harmonized hof's source code (just CUE here)
	harmonized := R.HarmonizeCue(base, source)

	// normal hof building process
	deps := R.FetchGoDeps(base, harmonized)
	builder := R.BuildHof(deps, harmonized)
	runner := R.HofImage(builder)

	// self tests?

	// make sure things dagger things happen
	// want to make sure our code works before loading registry
	runner.Sync(R.Ctx)
	out, err := runner.Stdout(R.Ctx)
	fmt.Println(out)
	checkErr(err)

	// load our registry
	reg, err := registry.Load(
		R.HofVer,
		R.CueVer,
		R.GoVer,
		R.ContainerRuntime,
		R.ContainerVersion,
	)
	checkErr(err)

	if R.Verbose {
		fmt.Printf("%# v", pretty.Formatter(reg))
	}

	// prep container runtime sidecar
	daemon, err := R.DockerDaemonContainer()
	checkErr(err)

	// attach dockerd to the tester containers
	runner, err = R.AttachDaemonAsService(runner, daemon)
	checkErr(err)
	builder, err = R.AttachDaemonAsService(builder, daemon)
	checkErr(err)

	for gkey, group := range *reg {
		// should we skip this group?
		if R.RunGroup != "" && !strings.HasPrefix(gkey, R.RunGroup) {
			continue
		}

		for ckey, C := range group {
			// should we skip this case?
			if R.RunCase != "" && !strings.HasPrefix(ckey, R.RunCase) {
				continue
			}

			// get source from git
			repo := R.GitSource(C.URL, C.Ref)

			// setup tester container and do any harmonization
			var tester *dagger.Container
			switch C.Type {
			case "cli":
				// setup the runtime container with case source
				tester = runner.
					Pipeline(fmt.Sprintf("%s/%s", gkey, ckey)).
					WithWorkdir("/harmony").
					WithDirectory("/harmony", repo)

			case "pkg":
				// we use the builder container, which already has
				// our harmonized hof source

				// go mod edit...
				repo = R.HarmonizeCue(builder, repo)
				repo = R.HarmonizeHof(builder, repo)
				tester = builder.
					Pipeline(fmt.Sprintf("%s/%s", gkey, ckey)).
					WithWorkdir("/harmony").
					WithDirectory("/harmony", repo)
					// run tidy on code base for them
					tester = tester.WithExec([]string{"go", "mod", "tidy"})
			}

			// run tidy on code base for them
			tester = tester.WithExec([]string{"hof", "mod", "tidy"})

			// cache bust
			tester = tester.WithEnvVariable("CACHEBUST", time.Now().String())

			for _, script := range C.Scripts {
				s := tester.WithExec([]string{"bash", "-c", script})
				s.Sync(R.Ctx)
				out, err := s.Stdout(R.Ctx)
				fmt.Println(out)
				checkErr(err)
			}

		}
	}

}
