package ast

type Scope map[string]ASTNode

type Scoped interface {
	DefineInScope(name string, node ASTNode) error
	LookupInScope(path []string) (ASTNode, error)
}
