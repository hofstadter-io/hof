package ast

type Field struct {
	// Parser filled
	BaseNode

	Key   *Token
	Value ASTNode

	// Phases filled
	Name string
}

type FieldType struct {
	// Parser filled
	BaseNode

	Value ASTNode

	// Phases filled
}

type FieldValue struct {
	// Parser filled
	BaseNode

	Value ASTNode

	// Phases filled
}

