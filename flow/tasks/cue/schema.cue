package hof

Format: {
	@task(cue.Format)
	$task: "cue.Format"

	value: _

	Package:        string | *""
	Raw:            bool | *false
	Final:          bool | *false
	Concrete:       bool | *true
	Definitions:    bool | *true
	Optional:       bool | *true
	Hidden:         bool | *true
	Attributes:     bool | *true
	Docs:           bool | *true
	InlineImports:  bool | *false
	ErrorsAsValues: bool | *false

	out: string
}
