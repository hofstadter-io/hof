package ast

import "fmt"

type ScriptError struct {
	Message string
	Node    Node
	Err     error
}

func NewScriptError(msg string, node Node, err error) *ScriptError {
	return &ScriptError { msg, node, err }
}

func (e *ScriptError) Error() string {
	src := ""
	lines := e.Node.Script().Lines[e.Node.DocLine(): e.Node.EndLine()]
	if len(lines) == 1 {
		src = lines[0]
	} else {
		src = strings.Join(lines, "\n  > ")
	}
	msg := `%s:%d:%d:%d: %s %v
	  > %s
	`
	return fmt.Sprintf(e.Path, e.Node.BegLine() + 1, e.Node.EndLine() + 1, e.Message, e.Err, src)
}
