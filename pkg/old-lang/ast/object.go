package ast

type Object struct {
	// Parser filled
	BaseNode

	Fields []*Field

	// Phases filled
	SimpleScopedNode
}

