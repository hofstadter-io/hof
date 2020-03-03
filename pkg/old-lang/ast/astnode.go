package ast

type ASTNode interface {
	GetParseInfo() *ParseInfo
	Parent() ASTNode

	Visit(func(ASTNode) (error)) error
}

