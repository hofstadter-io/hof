package datamodel

import "cuelang.org/go/cue"

type Common struct {
	Status  string // should probably be an const [ok,dirty,no history]
	Label   string // label from CUE
	Version string // timestamp

	Value   cue.Value // this objects value
	Other   cue.Value // the other value for diff
	Diff    cue.Value // diff from other (checkpoint)
	Subsume error
}

type History struct {
	Curr *Datamodel   // Top-level Datamodel, there should only be one history
	Prev *Datamodel   // The previous datamodel, when needed for comparison
	Past []*Datamodel // the full history of the datamodel
}

type Datamodel struct {
	Name    string
	Models  map[string]*Model
	Ordered []*Model
	History *History

	Common
}

type Model struct {
	Name   string
	Fields map[string]*Field

	Common
}

type Field struct {
	Name string
	Type string

	Common
}
