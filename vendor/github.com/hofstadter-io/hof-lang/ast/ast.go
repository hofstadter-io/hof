package ast

type ASTNode interface {
	ToData() (interface{}, error)
	FromData(interface{}) (ASTNode, error)

	String(indent string) (string, error)
	Print(indent string)
	Pretty(indent string)

	Visit(func(ASTNode, interface{}) (interface{}, error)) error
}

