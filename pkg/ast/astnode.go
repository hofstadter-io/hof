package ast

type ASTNode interface {
	GetParseInfo() *ParseInfo
	Visit(func(ASTNode) (error)) error
}

