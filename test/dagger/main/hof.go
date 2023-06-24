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

	err = R.HofVersion(tester)
	checkErr(err)

	err = R.TestCommandFmt(tester, source)
	checkErr(err)

	err = R.TestAdhocRender(tester, source)
	checkErr(err)

	return

	//err = R.SanityTest(c)
	//checkErr(err)

	//err = R.DockerTest(c)
	//checkErr(err)

	//d, err := R.BuildHofMatrix(c)
	//checkErr(err)

	//ok, err := d.Export(R.Ctx, ".")
	//checkErr(err)
	//if !ok {
	//  panic("unable to write matrix build outputs")
	//}

	//err = R.RenderTests(c)
	//checkErr(err)

}

