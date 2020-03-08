package ast

func (N HofFile) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N Definition) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N TypeDecl) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N Object) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N Field) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N Array) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N PathExpr) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N TokenPath) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N BracePath) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N RangeExpr) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N Token) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N Integer) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N Decimal) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}

func (N Bool) Visit(V func(ASTNode, interface{}) (interface{}, error)) error {

	return nil
}
