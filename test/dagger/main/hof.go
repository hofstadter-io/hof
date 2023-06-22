package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	hdagger "github.com/hofstadter-io/hof/test/dagger"
)

// so we don't have to pass these around everywhere
type Runtime struct {
	ctx    context.Context
	client *dagger.Client
}

type stage func (*dagger.Container) (*dagger.Container, error)

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

	var c *dagger.Container

	c, err = R.BuildBase(nil)
	checkErr(err)

	c, err = R.LocalCodeAndDeps(c)
	checkErr(err)

	c, err = R.BuildHof(c)
	checkErr(err)

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

	err = R.RenderTests(c)
	checkErr(err)

}

