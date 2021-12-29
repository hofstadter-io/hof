package datamodel

import "cuelang.org/go/cue"

type common struct {
	status string // should probably be an const [ok,dirty,no history]
	label  string
	value  cue.Value
}

type History struct {
	Curr *Datamodel
	Past []*Datamodel

	Other *Datamodel // set to a comparison version if available
	Diff  cue.Value  // Diff from other -> curr
}

type Datamodel struct {
	Name    string
	Models  map[string]*Model
	Ordered []*Model
	History *History

	version string    // timestamp
	diff    cue.Value // diff from last checkpoint
	common
}

type Model struct {
	Name   string
	Fields map[string]*Field

	common
}

type Field struct {
	Name string
	Type string

	common
}
