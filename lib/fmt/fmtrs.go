package fmt

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const ContainerPrefix = "hof-fmt-"

type Formatter struct {
	// name, same as tools/%
	Name  string

	// Info
	Running   bool
	Port      string
	Container *types.Container
}

var formatters map[string]Formatter

var fmtrNames = []string{
	"black",
	"prettier",
}

func init() {
	formatters = make(map[string]Formatter)
	for _,fmtr := range fmtrNames {
		formatters[fmtr] = Formatter{Name: fmtr}
	}
}

func updateFormatterStatus() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	// reset formatters
	for name, fmtr := range formatters {
		fmtr.Running = false
		fmtr.Container = nil
		formatters[name] = fmtr
	}

	for _, container := range containers {
		// extract name
		name := container.Names[0]
		name = strings.TrimPrefix(name, "/" + ContainerPrefix)

		// get fmtr
		fmtr := formatters[name]

		// always set running, otherwise it would not be in the lines
		fmtr.Running = true

		p := 10000000
		for _, port := range container.Ports {
			P := int(port.PublicPort)
			if P < p {
				p = P
			}
		}

		fmtr.Port = fmt.Sprint(p)

		// save container to fmtr
		fmtr.Container = &container

		formatters[name] = fmtr
	}

	return nil
}
