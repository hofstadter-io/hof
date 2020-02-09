package ast

type ASTNode interface {
	Visit(func(ASTNode) (error)) error
}
