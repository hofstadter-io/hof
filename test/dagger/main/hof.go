package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	hdagger "github.com/hofstadter-io/hof/test/dagger"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
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
	base := R.BaseContainer()
	code := R.WithCodeAndDeps(base, source)
	builder := R.BuildHof(code)
	runner := R.RuntimeContainer(builder)

	//
	// TESTS
	//

	tester := R.SetupTestingEnv(runner)
	tester = tester.Pipeline("TESTS")

	// bust cache before testing
	// tester = tester.WithEnvVariable("CACHE", time.Now().String())

	err = R.HofVersion(tester)
	checkErr(err)

	// so we don't fail fast
	errs := make(map[string]error)

	err = R.TestCommandFmt(tester, source)
	errs["fmt"] = err

	err = R.TestAdhocRender(tester, source)
	errs["render"] = err

	err = R.TestMod(tester, source)
	errs["mod"] = err

	err = R.TestCreate(tester, source)
	errs["create"] = err

	hadErr := false
	for key, err := range errs {
		if err != nil {
			fmt.Println("error:", key, err)
			hadErr = true
		}
	}
	if hadErr {
		os.Exit(1)
	}
}
