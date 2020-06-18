package ast

type CmdExpect int

const (
	None CmdExpect = iota
	Pass
	Fail
	Skip
)

type Cmd struct {
	NodeBase

	Exp  CmdExpect
	Cmd  string
	Args []string
	Bg   bool
}

func (P *Parser) parseCmd() error {

	return nil
}
