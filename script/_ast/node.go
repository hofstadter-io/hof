package ast

import "strings"

type Node interface {
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

type NodeBase struct {
	script  *Script
	docLine int
	begLine int
	endLine int
	name    string
	comment string
}

func (n *NodeBase) Script() *Script {
	return n.script
}

func (n *NodeBase) SetScript(s *Script) {
	n.script = s
}

func (n *NodeBase) DocLine() int {
	return n.docLine
}
func (n *NodeBase) SetDocLine(i int) {
	n.docLine = i
}

func (n *NodeBase) BegLine() int {
	return n.begLine
}
func (n *NodeBase) SetBegLine(i int) {
	n.begLine = i
}

func (n *NodeBase) EndLine() int {
	return n.endLine
}
func (n *NodeBase) SetEndLine(i int) {
	n.endLine = i
}

func (n *NodeBase) String() string {
	return strings.Join(n.script.Lines[n.docLine:n.endLine+1], "\n")
}

func (n *NodeBase) Name() string {
	return n.name
}

func (n *NodeBase) Comment() string {
	return n.comment
}
func (n *NodeBase) AddComment(c string) {
	n.comment += c
}
