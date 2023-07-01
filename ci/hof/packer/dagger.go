package main

import (
	"context"
	"fmt"
	"os"
	gouser "os/user"
	"time"

	"dagger.io/dagger"
	hdagger "github.com/hofstadter-io/hof/test/dagger"
)

var (
	// the os user running this pipeline
	// used for vm login & auth
	user string

	runtimes = []string{
		"docker",
		"nerdctl",
		"nerdctl-rootless",
		"podman",
	}
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	ctx := context.Background()

	u, err := gouser.Current()
	checkErr(err)
	user = u.Username

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
		Exclude: []string{"cue.mod/pkg", "docs", "next", ".git"},
	})

	//
	// Building Hof
	//
	base := R.GolangImage()
	deps := R.FetchDeps(base, source)
	builder := R.BuildHof(deps, source)
	hof := builder.File("hof")

	gcloud := R.GcloudImage()
	gcloud = R.WithLocalGcloudConfig(gcloud)

	for _, runtime := range runtimes {
		vmName := fmt.Sprintf("%s-fmt-test-%s", user, runtime)
		vmFamily := fmt.Sprintf("hof-debian-%s", runtime)
		t := gcloud.Pipeline(vmName)
		t = t.WithEnvVariable("CACHEBUST", time.Now().String())
		t = WithBootVM(t, vmName, vmFamily)

		// any runtime extra pre-steps before testing
		// we really want to test that it is permission issue and advise the user
		// we should also capture this as a test, so we need a setup where this fails intentionally
		switch runtime {
		case "docker":	
			// will probably need something like this for nerdctl too
			t = WithGcloudRemoteCommand(t, vmName, "sudo usermod -aG docker $USER")	

		case "nerdctl":
			// https://github.com/containerd/nerdctl/blob/main/docs/faq.md#does-nerdctl-have-an-equivalent-of-sudo-usermod--ag-docker-user-
			// make a user home bin and add to path
			t = WithGcloudRemoteCommand(t, vmName, "mkdir -p $HOME/bin && chmod 700 $HOME/bin && echo 'PATH=$HOME/bin:$PATH' >> .profile")	
			// copy nerdctl and set bits appropriatedly
			t = WithGcloudRemoteCommand(t, vmName, "cp /usr/local/bin/nerdctl $HOME/bin && sudo chown root $HOME/bin/nerdctl && sudo chmod +s $HOME/bin/nerdctl")	
		}

		// remote commands to run
		t = WithGcloudSendFile(t, vmName, "/usr/local/bin/hof", hof, true)
		t = WithGcloudRemoteCommand(t, vmName, "hof version")
		t = WithGcloudRemoteCommand(t, vmName, "hof fmt pull all@v0.6.8-rc.5")

		// sync to run them for real
		_, err = t.Sync(R.Ctx)
		hadErr := false
		if err != nil {
			hadErr = true
			fmt.Println("an error!:", err)
		}

		// always try deleting, we mostly ignore the error here (less likely, will also error if not exists)
		d := gcloud.Pipeline("DELETE " + vmName)
		d = d.WithEnvVariable("CACHEBUST", time.Now().String())
		d = WithDeleteVM(d, vmName)
		_, err := d.Sync(R.Ctx)
		if err != nil {
			fmt.Println("deleting error!:", err)
		}

		// stop if we had an error
		if hadErr {
			fmt.Println("stopping b/c error")
			break
		}
	}

}

func WithBootVM(gcloud *dagger.Container, name, imageFamily string) (*dagger.Container) {
	args := []string{
		"gcloud",
		"compute",
		"instances",
		"create",
		name,
		"--zone=us-central1-a",
		"--machine-type=n2-standard-2",
		"--image-family=" + imageFamily,
	}

	return gcloud.WithExec(args)
}

func WithDeleteVM(gcloud *dagger.Container, name string) (*dagger.Container) {
	args := []string{
		"gcloud",
		"compute",
		"instances",
		"delete",
		"--quiet",
		name,
		"--zone=us-central1-a",
	}

	gcloud = gcloud.WithExec(args)

	return gcloud
}

func WithGcloudSendFile(gcloud *dagger.Container, name, remotePath string, file *dagger.File, sudo bool) (*dagger.Container) {
	tmpPath := "/file-to-copy"
	// add file in container
	c := gcloud.WithFile(tmpPath, file)

	// send file from container
	c = c.WithExec([]string{
		"gcloud",
		"compute",
		"scp",
		"--zone=us-central1-a",
		tmpPath,
		user + "@" + name + ":file-copied-tmp",
	})

	// build up remote copy
	mv := []string{
		"gcloud",
		"compute",
		"ssh",
		"--zone=us-central1-a",
		user + "@" + name,
		"--",
	}
	if sudo {
		mv = append(mv, "sudo")
	}
	mv = append(mv, "mv", "file-copied-tmp", remotePath)

	c = c.WithExec(mv)

	return c
}

func WithGcloudScp(gcloud *dagger.Container, name string, dir *dagger.Directory) (*dagger.Container) {

	c := gcloud.WithDirectory("/src", dir)

	// tar file

	// copy tar
	c = c.WithExec([]string{
		"gcloud",
		"compute",
		"scp",
		"--recurse",
		"--zone=us-central1-a",
		"/src",
		user + "@" + name + ":src",
	})

	// untar

	return c
}

func WithGcloudRemoteCommand(gcloud *dagger.Container, name string, cmd string) (*dagger.Container) {
	return gcloud.WithExec([]string{
		"gcloud",
		"compute",
		"ssh",
		user + "@" + name,
		"--zone=us-central1-a",
		"--",
		"bash", 
		"-c",
		fmt.Sprintf("'set -euo pipefail; %s'", cmd),
	})
}
