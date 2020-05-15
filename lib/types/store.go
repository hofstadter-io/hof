package types

import (
	"cuelang.org/go/cue"
)

type Store struct {
	ID      string
	Type    string
	Version string

	CueValue  cue.Value
}
