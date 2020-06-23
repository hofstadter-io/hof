package ast

type Node interface {
	Clone()         Node
	CloneNodeBase() NodeBase

	Script()  *Script

	DocLine() int
	SetDocLine(int)
	BegLine() int
	SetBegLine(int)
	EndLine() int
	SetEndLine(int)

	String()  string
	Name()    string
	Comment() string
	AddComment(string)
}
