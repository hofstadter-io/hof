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
		"none",
		//"docker",
		//"nerdctl",
		//"nerdctl-rootless",
		//"podman",
	}
	arches = []string{
		"amd",
		"arm",
	}

	machSize = "standard-2"
	diskSize = "200GB"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
 *
 *  Note, this script must be run from the repo root
 *
 */

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

	// gcloud image to run commands in
	gcloud := R.GcloudImage()

	// mount local config & secrets
	gcloud = R.WithLocalGcloudConfig(gcloud)
	// gcloud = R.WithLocalSSHDir(gcloud)

	//
	// Testing on matrix of {arch}x{runtime}
	//
	hadErr := false
	for _, arch := range arches {
		// build hof the normal way
		base := R.GolangImage(fmt.Sprintf("linux/%s64", arch))
		deps := R.FetchDeps(base, source)
		builder := R.BuildHof(deps, source)
		hof := builder.File("hof")

		for _, runtime := range runtimes {
			vmName := fmt.Sprintf("%s-fmt-test-%s-%s", user, runtime, arch)
			t := gcloud.Pipeline(vmName)
			t = t.WithEnvVariable("CACHEBUST", time.Now().String())
			
			// start VM
			vmFamily := fmt.Sprintf("debian-%s-%s", runtime, arch)
			t = WithBootVM(t, vmName, vmFamily, arch)

			// any runtime extra pre-steps before testing
			// we really want to test that it is permission issue and advise the user
			// we should also capture this as a test, so we need a setup where this fails intentionally
			switch runtime {
			case "docker":	
				// will probably need something like this for nerdctl too
				t = WithGcloudRemoteBash(t, vmName, "sudo usermod -aG docker $USER")	

			case "nerdctl":
				t = WithGcloudRemoteBash(t, vmName, "sudo nerdctl apparmor load")	
				// https://github.com/containerd/nerdctl/blob/main/docs/faq.md#does-nerdctl-have-an-equivalent-of-sudo-usermod--ag-docker-user-
				// make a user home bin and add to path
				t = WithGcloudRemoteBash(t, vmName, "mkdir -p $HOME/bin && chmod 700 $HOME/bin && echo 'PATH=$HOME/bin:$PATH' >> .profile")	
				// copy nerdctl and set bits appropriatedly
				t = WithGcloudRemoteBash(t, vmName, "cp /usr/local/bin/nerdctl $HOME/bin && sudo chown $(id -u):$(id -g) $HOME/bin/nerdctl && sudo chmod 0755 $HOME/bin/nerdctl && sudo chown root $HOME/bin/nerdctl && sudo chmod +s $HOME/bin/nerdctl")	
				t = WithGcloudRemoteBash(t, vmName, "nerdctl version")	

			case "nerdctl-rootless":
				// ensure the current user can run nerdctl
				t = WithGcloudRemoteBash(t, vmName, "containerd-rootless-setuptool.sh install")	
				t = WithGcloudRemoteBash(t, vmName, "nerdctl version")	
			}

			// remote commands to run
			t = WithGcloudSendFile(t, vmName, "/usr/local/bin/hof", hof, true)
			t = WithGcloudRemoteBash(t, vmName, "hof version")
			t = WithGcloudRemoteBash(t, vmName, "hof fmt pull prettier@v0.6.8-rc.5")
			t = WithGcloudRemoteBash(t, vmName, "hof fmt start prettier@v0.6.8-rc.5")
			t = WithGcloudRemoteBash(t, vmName, "hof fmt status")
			t = WithGcloudRemoteBash(t, vmName, "hof fmt test prettier")
			t = WithGcloudRemoteBash(t, vmName, "hof fmt stop")

			// sync to run them for real
			_, err = t.Sync(R.Ctx)
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
				break
			}
		} // end runtime loop
		// stop if we had an error
		if hadErr {
			fmt.Println("stopping b/c error")
			break
		}
	} // end arch loop
} // end main

func WithBootVM(gcloud *dagger.Container, name, imageFamily, arch string) (*dagger.Container) {
	mach := "n2-" + machSize
	if arch == "arm" {
		mach = "t2a-" + machSize
	}
	args := []string{
		"gcloud",
		"compute",
		"instances",
		"create",
		name,
		"--zone=us-central1-a",
		"--machine-type=" + mach,
		"--boot-disk-size=" + diskSize,
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

func WithGcloudRemoteBash(gcloud *dagger.Container, name string, cmd string) (*dagger.Container) {
	return gcloud.WithExec([]string{
		"gcloud",
		"compute",
		"ssh",
		user + "@" + name,
		"--zone=us-central1-a",
		"--",
		"bash", 
		"--login",
		"-c",
		fmt.Sprintf("'set -euo pipefail; %s'", cmd),
	})
}
