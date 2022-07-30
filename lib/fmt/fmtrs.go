package fmt

import (
	"github.com/docker/docker/api/types"

	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
)

const ContainerPrefix = "hof-fmt-"

var defaultVersion = "dirty"

func init() {
	v := verinfo.Version
	if v != "Local" {
		defaultVersion = v
	}
}


type Formatter struct {
	// name, same as tools/%
	Name    string
	Version string
	Available []string

	// Info
	Running   bool
	Port      string
	Container *types.Container
	Images    []*types.ImageSummary

	Config  interface{}
	Default interface{}
}

var formatters map[string]*Formatter

func init() {
	formatters = make(map[string]*Formatter)
	for _,fmtr := range fmtrNames {
		d := fmtrDefaults[fmtr]
		formatters[fmtr] = &Formatter{Name: fmtr, Default: d, Version: defaultVersion}
	}
}

var fmtrNames = []string{
	"black",
	"prettier",
}

var fmtrDefaults = map[string]interface{}{
}
