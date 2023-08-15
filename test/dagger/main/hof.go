package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"dagger.io/dagger"
	"github.com/spf13/pflag"

	hdagger "github.com/hofstadter-io/hof/test/dagger"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var RUNTIME string
var TESTS string

func init() {
	pflag.StringVarP(&RUNTIME, "runtime", "r", "docker", "container runtime to use [docker, nerdctl, podman, none]")
	pflag.StringVarP(&TESTS, "tests", "t", "", "tests to run, comma separated")
}

func main() {
	pflag.Parse()

	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	R := &hdagger.Runtime{
		Ctx:    ctx,
		Client: client,
	}

	// load hof's code from the host
	// todo, find repo root with git
	source := R.Client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{"cue.mod/pkg", "docs", "next"},
	})

	//
	// Building
	//
	base := R.GolangImage("linux/amd64")
	deps := R.FetchDeps(base, source)
	builder := R.BuildHof(deps, source)
	runner := R.RuntimeContainer(builder, "linux/amd64")

	//
	// TESTS
	//

	tester := R.SetupTestingEnv(runner, source)
	tester = tester.Pipeline("TESTS")

	switch RUNTIME {
	case "none":
		tester = tester.WithEnvVariable("HOF_CONTAINER_RUNTIME", "none")
	
	case "docker" :
		// add tools
		tester = R.AddDockerCLI(tester)

		// attach dockerd to the tester container
		daemon, err := R.DockerDaemonContainer()
		check(err)
		tester, err = R.AttachDaemonAsService(tester, daemon)
		check(err)
	}

	// bust cache before testing
	tester = tester.
		WithEnvVariable("CACHEBUST", time.Now().String()).
		WithExec([]string{"env"})

	// sync the graph
	tester, err = tester.Sync(R.Ctx)
	if err != nil {
		check(err)
	}

	// run hof version as a first sanity test
	err = R.HofVersion(tester)
	check(err)

	// build up a test map so we can easily select which ones to run
	tests := make(map[string]func() error)
	tests["render"] = func() error {
		return R.TestAdhocRender(tester, source)
	}
	tests["create"] = func() error {
		return R.TestCreate(tester, source)
	}
	tests["flow"] = func() error {
		return R.TestFlow(tester, source)
	}
	tests["st"] = func() error {
		return R.TestStructural(tester, source)
	}
	tests["dm"] = func() error {
		return R.TestDatamodel(tester, source)
	}
	tests["mod"] = func() error {
		return  R.TestMod(tester, source)
	}
	tests["fmt"] = func() error {
		return R.TestCommandFmt(tester, source)
	}

	// decide what tests to run
	ts := []string{}
	if TESTS == "" {
		for k := range tests {
			ts = append(ts,k)
		}
	} else {
		ts = strings.Split(TESTS,",")
	}

	// run tests
	for _,t := range ts {
		fn, ok := tests[t]
		if !ok {
			fmt.Println("unknown test %q", t)
			os.Exit(1)
		}

		err := fn()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	tester.WithExec([]string{"echo", "finished!"})
	fmt.Println("tests finished!")
}
