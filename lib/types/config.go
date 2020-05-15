package types

import (
	"cuelang.org/go/cue"
)

type Config struct {
	Modelsets map[string]Modelset
	Stores    map[string]Store

	CueValue  cue.Value
}
