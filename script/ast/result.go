package ast

import (
	"io"
	"strings"
	"time"
)

type Result struct {
	Node Node

	BegTime time.Time
	EndTime time.Time

	StdinR io.ReadCloser
	StdinW io.Writer
	Stdout io.Writer
	Stderr io.Writer
	Status int
	Errors []error

	Parent *Result
	Nodes  []*Result
}

func NewResult(n Node, p *Result) *Result {
	r := &Result {
		Node: n,
		Parent: p,
		Status: -42,
	}
	r.DefaultWriters()
	n.SetResult(r)
	return r
}

func (R *Result) AddError(err error) {
	R.Errors = append(R.Errors, err)
}

func (R *Result) AddResult(r *Result) {
	R.Nodes = append(R.Nodes, r)
}

func (R *Result) DefaultWriters() {
	R.Stdout = new(strings.Builder)
	R.Stderr = new(strings.Builder)
}
