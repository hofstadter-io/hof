package ast

type ArrayDef struct {
	// Parser filled
	BaseNode

	Path      *TokenPath

	// Phases filled
}

type Array struct {
	// Parser filled
	BaseNode

	Elems []ASTNode

	// Phases filled
}
