package types

import (
	"cuelang.org/go/cue"
)

type Secret struct {
	StoreName  string
	CredValues map[string]string

	CueValue cue.Value
}
