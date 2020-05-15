package types

import (
	"cuelang.org/go/cue"
)

type Creds struct {
	StoreName string
	CredValues map[string]string

	CueValue  cue.Value
}
