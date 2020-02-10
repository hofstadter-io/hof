package ast

import "fmt"

type ParseInfo struct {
	File   string
	Line   int
	Col    int
	Offset int
	Text   string
}


func (p ParseInfo) String() string {
	return fmt.Sprintf("%s %d:%d", p.File, p.Line, p.Col)
}

