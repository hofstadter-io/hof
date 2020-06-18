package ast

type Phase struct {
	NodeBase

	Phases []*Phase
	Nodes  []*Node
}

func (P *Parser) parsePhase() error {

	return nil
}
