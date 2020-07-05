package ast

type Node interface {
	Clone()         Node
	CloneNodeBase() NodeBase

	Script()  *Script
	SetScript(*Script)
	Result() *Result
	SetResult(*Result)

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
