package ast

type Token struct {
	// Parser filled
	BaseNode

	Value string
}

type TokenPath struct {
	// Parser filled
	BaseNode

	Paths []string

	// Phases filled
}

