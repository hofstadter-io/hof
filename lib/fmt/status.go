package fmt

import (
	"fmt"
	"strings"
	"time"

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

		c := container
		updateFmtrContainer(fmtr, &c)

		formatters[name] = fmtr
	}

	return nil
}

func getAndUpdateFmtrContainer(fmtr *Formatter) error {
	containers, err := container.GetContainers(ContainerPrefix + fmtr.Name)
	if err != nil {
		return err
	}

	for _, c := range containers {
		// extract name
		name := c.Names[0]
		name = strings.TrimPrefix(name, ContainerPrefix)
		if name == fmtr.Name {
			C := c
			updateFmtrContainer(fmtr, &C)
			break
		}
	}

	return nil
}

func updateFmtrContainer(fmtr *Formatter, container *container.Container) {

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
		fmtr.Container = container
}

func (fmtr *Formatter) WaitForRunning(retry int, delay time.Duration) error {

	// We should probably rethink how this works, such that
	// we minimize our exec out to docker (et al)
	// we can exec once and then check on all formatters

	// fmt.Println("wait-running.0:", fmtr.Name, fmtr.Status, fmtr.Running, fmtr.Ready)
	// return if already running
	if fmtr.Running {
		return nil
	}

	for i := 0; i < retry; i++ {

		err := getAndUpdateFmtrContainer(fmtr)
		if err != nil {
			return err
		}

		if fmtr.Running {
			return nil
		}

		time.Sleep(delay)
	}

	return fmt.Errorf("formatter %s never started", fmtr.Name)
}

func (fmtr *Formatter) WaitForReady(retry int, delay time.Duration) error {
	// fmt.Println("wait-ready.0:", fmtr.Name, fmtr.Status, fmtr.Running, fmtr.Ready)
	err := getAndUpdateFmtrContainer(fmtr)
	if err != nil {
		return err
	}

	// return if already ready
	if fmtr.Ready {
		return nil
	}

	// return error if not running
	if !fmtr.Running {
		return fmt.Errorf("formatter %s is not running", fmtr.Name)
	}

	// get ready check payload
	p, ok := fmtrReady[fmtr.Name]
	if !ok {
		fmt.Printf("warn: formatter %s does not have a ready config\n", fmtr.Name)
		return nil
	}

	payload := p.(map[string]any)

	for i := 0; i < retry; i++ {
		_, err := fmtr.Call("ready-check", []byte(payload["source"].(string)), payload["config"])
		// fmt.Println("wait-ready:", i, fmtr.Name, fmtr.Ready, err)
		// if no error, then ready
		if err == nil {
			fmtr.Ready = true
			return nil
		}
		time.Sleep(delay)
	}

	return fmt.Errorf("formatter %s is not ready", fmtr.Name)
}
