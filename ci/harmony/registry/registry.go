package registry

import (
	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/cuetils"
)

type Registry map[string]Group

type Group map[string]Case

type Case struct {
	Group string `json:"group"`
	Name  string `json:"name"`

	URL  string `json:"url"`
	Ref  string `json:"ref"`
	Type string	`json:"type"`

	Scripts []string `json:"scripts"`
}

func Load(hofver, cuever, gover, cruntime, cversion string) (*Registry, error) {

	// load ci/harmony/registry/*cue
	crt, err := cuetils.CueRuntimeFromEntrypoints([]string{"./ci/harmony/registry/"})
	if err != nil {
		return nil, err
	}

	// inject version details
	val := crt.CueValue
	val = val.FillPath(cue.ParsePath("versions"), map[string]string{
		"hof": hofver,
		"cue": cuever,
		"go": gover,
		"runtime": cruntime,
		"version": cversion,
	})

	reg := val.LookupPath(cue.ParsePath("Registry"))

	// decode into go
	var R Registry
	err = reg.Decode(&R)

	return &R, nil
}
