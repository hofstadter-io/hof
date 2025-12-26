package dag

import "cuelang.org/go/cue"

type StepKind struct {
	Kind  string `json:"$kind"`
	Value cue.Value
}

type Step map[string]any
