package main

import (
	"fmt"

	"github.com/hofstadter-io/hof/pkg/lang/hof/parser"
	"github.com/hofstadter-io/hof/pkg/lang/hof/token"
)

func main() {

	p, err := parser.AParser("test.cue")
	if err != nil {
		panic(err)
	}

	t := p.Tok()
	for t != token.EOF {
		fmt.Printf("%s %s\n", t.String(), p.Lit())
		p.Next()
		t = p.Tok()
	}
}
