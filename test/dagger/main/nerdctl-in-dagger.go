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
	checkErr(err)
	defer client.Close()

	R := &hdagger.Runtime{
		Ctx:    ctx,
		Client: client,
	}

	// load source code from the host, technically current directory
  // you can "go run ./path/to/dagger.go" from the repo root to get everything
	// you can also find repo root with git and load everything based on that
	//source := R.Client.Host().Directory(".", dagger.HostDirectoryOpts{
	//  Exclude: []string{"cue.mod/pkg", "docs", "next"},
	//})

	//
	// Building
	//

	c := R.Client.Container().From("debian:12")
	r := R.AddNerdctl(c)

	opts := dagger.ContainerWithExecOpts{ SkipEntrypoint: true, InsecureRootCapabilities: true }
	daemon := r.WithExec([]string{"containerd", "-address", "0.0.0.0:2375"}, opts)

	t := r.Pipeline("nerdctl/test")

	// t = t.WithEnvVariable("DOCKER_HOST", "tcp://global-containerd:2375")
	t = t.WithServiceBinding("global-containerd", daemon)
	t = withRuntimeTest("nerdctl", t)

	out, err := t.Stdout(R.Ctx)
	fmt.Println(out)
	checkErr(err)

}

const socat = `
set -euo pipefail

touch /run/containerd/containerd
socat -d -d UNIX-LISTEN:/run/containerd/containerd.sock,reuseaddr,mode=777 TCP:global-containerd:2375 &
sleep 1
ls -lh /run/containerd
%s
pkill socat
`

func withRuntimeTest(runtime string, c *dagger.Container) (*dagger.Container) {
	t := c.Pipeline(runtime + "/test")

	opts := dagger.ContainerWithExecOpts{ SkipEntrypoint: true, InsecureRootCapabilities: true }

	t = t.WithExec([]string{"mkdir", "/run/containerd"}, opts)
	t = t.WithExec([]string{"apt-get", "install", "-y", "socat"}, opts)

	t = t.WithEnvVariable("CACHE", time.Now().String())

	//t = t.WithEnvVariable("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	//t = t.WithEntrypoint([]string{})

	t = t.WithExec([]string{"bash", "-c", fmt.Sprintf(socat, runtime + " --version")}, opts)
	t = t.WithExec([]string{"bash", "-c", fmt.Sprintf(socat, runtime + " info")}, opts)

	t = t.WithExec([]string{"bash", "-c", fmt.Sprintf(socat, runtime + " pull hello-world")}, opts)
	t = t.WithExec([]string{"bash", "-c", fmt.Sprintf(socat, runtime + " run hello-world")}, opts)

	//t = t.WithExec([]string{bin, "pull", "hello-world"})
	//t = t.WithExec([]string{bin, "run", "hello-world"})

	//t = t.WithExec([]string{bin, "pull", "nginxdemos/hello"})
	//t = t.WithExec([]string{bin, "images"})
	//t = t.WithExec([]string{bin, "run", "-p", "4000:80", "-d", "nginxdemos/hello"})
	//t = t.WithExec([]string{"curl", "localhost:4000"})

	return t
}
