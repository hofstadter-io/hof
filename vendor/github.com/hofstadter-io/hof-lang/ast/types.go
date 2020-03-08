package ast

type ParseInfo struct {
	Line int
	Col int
	Text string
}

type HofFile struct {
	Name string
	Path string

	Definitions []Definition
}

type Definition struct {
	ParseInfo ParseInfo

	Name Token
	DSL  Token

	Body []ASTNode
}

type TypeDecl struct {
	Name Token
	Type Token
	Extra *Object
}

type Object struct {
	Fields []Field
}

type Field struct {
	Key   Token
	Value ASTNode
}

type Array struct {
	Elems []ASTNode
}





type PathExpr struct {
	PathList []ASTNode
}

type TokenPath struct {
	Value string
}

type RangeExpr struct {
	Low   int
	High  int
	Range bool
}

type BracePath struct {
	// Exprs []Expr
	Exprs []ASTNode
}

type Token struct {
	Value string
}

type Integer struct {
	Value int
}

type Decimal struct {
	Value float64
}

type Bool struct {
	Value bool
}
