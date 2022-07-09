package datamodel

import "cuelang.org/go/cue"

type Common struct {
	// User set
	Name   string
	Labels map[string]string

	// internally set
	Status    string // should probably be an const [ok,dirty,no history]
	Label     string // label from CUE
	Timestamp string // timestamp
	Version   string // @dm_ver()
	// TODO(tags), maybe labels too?

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
	Models  map[string]*Model

	// filled fields
	OrderedModels []*Model
	History *History

	Common
}

type Model struct {
	Fields map[string]*Field

	// Filled fields
	OrderedFields []*Field
  Path []string

	Common
}

type Field struct {
	Type string

	Common
}
