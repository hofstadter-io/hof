package datamodel

import "cuelang.org/go/cue"

type common struct {
	status string
	label  string
	value  cue.Value
}

type History struct {
	Curr *Datamodel
	Past []*Datamodel
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
