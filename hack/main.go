// This file is used for hacking on functions in the repository

package main

import (
	"fmt"
	"os"

	// "github.com/kr/pretty"

	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/hof"
)

var HNS []*hof.Node

func checkErr(err error) {
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}

func main() {
	args := os.Args[1:]
	fmt.Println("hack!", args)

	CRT, err := cuetils.CueRuntimeFromEntrypoints(args)
	checkErr(err)

	HNS, err = hof.FindHofs(CRT.CueValue)
	checkErr(err)


	fmt.Println("found:")
	for _, hn := range HNS {
		printNodeTree(hn, "  ", "")
	}
}

func printNodeTree(hn *hof.Node, spacer, indent string) {
	s := fmt.Sprintf("%s%s - ", indent, hn.Hof.Path)
	dm := hn.Hof.Datamodel
	if dm.Root {
		s += " @datamodel()"
	}
	if dm.Node {
		s += " @node()"
	}
	if dm.History {
		s += " @history()"
	}
	if dm.Ordered {
		s += " @ordered()"
	}
	if dm.Cue {
		s += " @cue()"
	}
	
	// print line
	fmt.Println(s)
	// pretty.Println(hn.Hof)

	// recurse
	for _, c := range hn.Children {
		printNodeTree(c, spacer, indent+spacer)
	}
}
