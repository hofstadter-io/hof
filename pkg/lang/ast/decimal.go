package ast

type DecimalDef struct {
	// Parser filled
	BaseNode
}

type Decimal struct {
	// Parser filled
	BaseNode

	Value float64
}
