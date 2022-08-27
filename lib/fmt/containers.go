package fmt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var dockerCli *client.Client

func initDockerCli() (err error) {
	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error: hof fmt requires docker\n%w", err)
	}
	return nil
}

func updateFormatterStatus() error {
	iFilter := filters.NewArgs(filters.Arg("reference", "hofstadter/fmt-*"))
	images, err := dockerCli.ImageList(
		context.Background(),
		types.ImageListOptions{
			Filters: iFilter,
		},
	)
	if err != nil {
		return err
	}

	cFilter := filters.NewArgs(filters.Arg("name", "hof-fmt-"))
	containers, err := dockerCli.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			All: true,
			Filters: cFilter,
		},
	)
	if err != nil {
		return err
	}

	// reset formatters
	for _, fmtr := range formatters {
		fmtr.Running = false
		fmtr.Container = nil
		fmtr.Available = make([]string, 0)
	}

	for _, image := range images {
		added := false
		for _, tag := range image.RepoTags {
			parts := strings.Split(tag, ":")
			repo, ver := parts[0], parts[1]
			name := strings.TrimPrefix(repo, "hofstadter/fmt-")
			fmtr := formatters[name]
			fmtr.Available = append(fmtr.Available, ver)
			if !added {
				fmtr.Images = append(fmtr.Images, &image)
				added = true
			}
		}
	}

	for _, container := range containers {
		// extract name
		name := container.Names[0]
		name = strings.TrimPrefix(name, "/" + ContainerPrefix)

		// get fmtr
		fmtr := formatters[name]

		// always set running, otherwise it would not be in the lines
		fmtr.Running = true

		p := 100000
		for _, port := range container.Ports {
			P := int(port.PublicPort)
			if P < p {
				p = P
			}
		}

		if p != 100000 {
			fmtr.Port = fmt.Sprint(p)
		}

		// save container to fmtr
		c := container
		fmtr.Container = &c

		formatters[name] = fmtr
	}

	return nil
}

func startContainer(fmtr string) error {
	// just try to pull, if already present this will not be noticed
	err := maybePullContainer(fmtr)
	if err != nil {
		return err
	}

	// maybe stop here, by ignoring error
	// we do this to cleanup and stopped / done images
	// (typical of docker when a user restarts their computer)
	stopContainer(fmtr)

	// now we can safely (re)create our container
	ret, err := dockerCli.ContainerCreate(
		context.Background(),

		// config
		// todo, maybe walk back versions? (dirty -> latest release)
		&container.Config{
			Image: fmt.Sprintf("hofstadter/fmt-%s:%s", fmtr, defaultVersion),
		},

		// hostConfig
		&container.HostConfig{
			PublishAllPorts: true,
		},

		// netConfig
		nil,

		// todo, need to consider mac arm here?
		// platform
		nil,

		// name
		fmt.Sprintf("hof-fmt-%s", fmtr),
	)

	if err != nil {
		return err
	}

	err = dockerCli.ContainerStart(
		context.Background(),
		ret.ID,
		types.ContainerStartOptions{},
	)

	// TODO, add alive command and wait for ready
	time.Sleep(500*time.Millisecond)

	return err
}

func stopContainer(fmtr string) error {
	return dockerCli.ContainerRemove(
		context.Background(),
		fmt.Sprintf("hof-fmt-%s", fmtr),
		types.ContainerRemoveOptions{ Force: true },
	)
}

func pullContainer(fmtr string) error {
	if defaultVersion == "dirty" {
		return fmt.Errorf("%s: You have local changes to hof, run 'make formatters' instead", fmtr)
	}
	ref := fmt.Sprintf("hofstadter/fmt-%s:%s", fmtr, defaultVersion)
	fmt.Println("pulling:", ref)

	r, err := dockerCli.ImagePull(context.Background(), ref, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	_, err = dockerCli.ImageLoad(context.Background(), r, false)
	if err != nil {
		return err
	}

	return nil
}

func maybePullContainer(fmtr string) error {
	if defaultVersion == "dirty" {
		return nil
	}
	F := formatters[fmtr]

	found := false
	for _, a := range F.Available {
		if a == fmtr {
			found = true
			break
		}
	}

	if !found {
		return pullContainer(fmtr)
	}

	return nil
}
