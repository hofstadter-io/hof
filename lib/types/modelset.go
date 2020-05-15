package types

import (
	"cuelang.org/go/cue"
)

type Modelset struct {
	Name string

	Entry  string
	CueValue  cue.Value

	Models map[string]Model
}

