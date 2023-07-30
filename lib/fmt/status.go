package fmt

import (
	"fmt"
	"strings"
	"github.com/hofstadter-io/hof/lib/container"
)

func UpdateFormatterStatus() error {
	images, err := container.GetImages(fmt.Sprintf("%s/fmt-", CONTAINER_REPO))
	if err != nil {
		return fmt.Errorf("get images: %w", err)
	}
	containers, err := container.GetContainers(ContainerPrefix)
	if err != nil {
		return fmt.Errorf("get containers: %w", err)
	}

	// reset formatters
	for _, fmtr := range formatters {
		fmtr.Running = false
		fmtr.Container = nil
		fmtr.Available = make([]string, 0)
		fmtr.Images = []*container.Image{}
	}

	for _, image := range images {
		img := image
		name := strings.TrimPrefix(image.Repository, fmt.Sprintf("%s/fmt-", CONTAINER_REPO))
		fmtr := formatters[name]
		if fmtr == nil {
			fmt.Printf("%q %# +v\n", name, image)
			continue
		}
		if len(image.RepoTags) > 0 {
			fmtr.Available = append(fmtr.Available, image.RepoTags...)
		} else {
			// podman...?
			fmtr.Available = append(fmtr.Available, image.RepoTags...)
		}
		fmtr.Images = append(fmtr.Images, &img)

		// fmt.Println(name, fmtr, image)
	}

	for _, container := range containers {
		// extract name
		name := container.Names[0]
		name = strings.TrimPrefix(name, ContainerPrefix)

		// get fmtr
		fmtr := formatters[name]

		fmtr.Status = container.State

		// determine the container status
		// TODO, we have various places where we ETL the container runtime responses
		//       we should move this to CUE, centralize it, and support multiple versions
		if container.State == "running" {
			fmtr.Running = true
		} else {
			fmtr.Running = false
		}

		p := 100000
		for _, port := range container.Ports {
			if port < p {
				p = port
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

