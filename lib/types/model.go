package types

import (
	"cuelang.org/go/cue"
)

type Model struct {
	Name string

	CueValue cue.Value
}
