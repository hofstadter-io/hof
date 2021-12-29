package datamodel

type History struct {
	Curr *Datamodel
	Past []*Datamodel
}

type Datamodel struct {
	Name    string
	Version string
	Models  []*Model
	History *History

	status string
}

type Model struct {
	Name   string
	Fields map[string]*Field
}

type Field struct {
	Name string
	Type string
}
