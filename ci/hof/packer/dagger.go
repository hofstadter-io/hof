package main

import (
	"context"
	"fmt"
	"os"
	"time"

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
	// Building Hof
	//
	//base := R.GolangImage()
	//deps := R.FetchDeps(base, source)
	//builder := R.BuildHof(deps, source)
	//hof := builder.File("hof")

	gcloud := R.GcloudImage()

	name := "test-debian"
	t := gcloud.WithEnvVariable("CACHEBUST", time.Now().String())
	t = WithBootVM(t, name)
	t = WithGcloudScp(t, name, source)
	t = WithGcloudRemoteCommand(t, name, "ls -l /src")
	t = WithDeleteVM(t, name)

	t.Sync(R.Ctx)
}

func WithBootVM(gcloud *dagger.Container, name string) (*dagger.Container) {
	args := []string{
		"gcloud",
		"compute",
		"instances",
		"create",
		name,
		"--zone=us-central1-a",
		"--machine-type=n2-standard-2",
		"--image-family=hof-debian-nerdctl",
	}

	return gcloud.WithExec(args)
}

func WithDeleteVM(gcloud *dagger.Container, name string) (*dagger.Container) {
	args := []string{
		"gcloud",
		"compute",
		"instances",
		"delete",
		name,
		"--zone=us-central1-a",
	}

	return gcloud.WithExec(args)
}

func WithGcloudScp(gcloud *dagger.Container, name string, dir *dagger.Directory) (*dagger.Container) {

	c := gcloud.WithDirectory("/src", dir)
	c = c.WithExec([]string{
		"gcloud",
		"compute",
		"scp",
		"--recurse",
		"--zone=us-central1-a",
		"/src",
		name + ":src",
	})

	return c
}

func WithGcloudRemoteCommand(gcloud *dagger.Container, name string, cmd string) (*dagger.Container) {
	return gcloud.WithExec([]string{
		"gcloud",
		"compute",
		"ssh",
		name,
		"--zone=us-central1-a",
		"--",
		"bash", 
		"-c",
		cmd,
	})
}
