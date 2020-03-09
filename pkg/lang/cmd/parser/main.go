package main

import (
	"fmt"

	"github.com/hofstadter-io/hof/pkg/lang/hof/parser"
)

func main() {

	f, err := parser.ParseFile("test.cue", nil)
	if err != nil {
		panic(err)
	}
	for _, d := range f.Decls {
		fmt.Printf("%s\n", parser.DebugStr(d))
	}
}
