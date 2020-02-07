package ast

type ASTNode interface {
	Walk(func(ASTNode) (error)) error
}
