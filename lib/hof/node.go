package hof

import (
	"cuelang.org/go/cue"
)

type Node[T any] struct {
	Hof Hof

	// containing value
	Value    cue.Value

	// The wrapping type
	T *T

	// heirarchy of tracked values
	Parent   *Node[T]
	// we (this node) are in between
	Children []*Node[T]

	// cue paths to get up/down hierarchy
}

func New[T any](label string, val cue.Value, curr *T, parent *Node[T]) *Node[T] {
	n := &Node[T]{
		Hof: Hof{
			Path: val.Path().String(),
			Label: label,
		},
		Value:    val,
		T:        curr,
		Parent:   parent,
		Children: make([]*Node[T], 0),
	}

	return n
}

func Upgrade[S, T any](src *Node[S], t func(*Node[T])*T, parent *Node[T]) *Node[T] {
	n := &Node[T]{
		Hof:      src.Hof,
		Value:    src.Value,
		Parent:   parent,
		Children: make([]*Node[T], 0, len(src.Children)),
	}

	n.T = t(n)

	// walk, upgrading children
	for _, c := range src.Children {
		u := Upgrade[S,T](c, t, n)
		n.Children = append(n.Children, u)
	}

	return n
}
